package coordinator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gaganpreet-S1ngh/dts-proto/pb"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type workerInfo struct {
	serviceClient   pb.WorkerServiceClient
	Connection      *grpc.ClientConn
	heartBeatMisses uint8
	address         string
}

type Service interface {
	RegisterWorker(workerID uint64, address string) error
	ManageWorkerPool(ctx context.Context)
	ExecuteAllScheduledTasks() error
	ScanDatabase(ctx context.Context)
	UpdateTaskStatus(ctx context.Context, taskID uint64, status pb.TaskStatus, startedAt int64, completedAt int64, failedAt int64) error
}

type service struct {
	repo                Repository
	workerPoolMutex     sync.Mutex
	workerPoolKeysMutex sync.Mutex
	workerPool          map[uint64]*workerInfo
	workerPoolKeys      []uint64
	roundRobinIndex     uint32
	wg                  sync.WaitGroup
}

func NewService(repo Repository) Service {
	return &service{
		repo:            repo,
		workerPool:      make(map[uint64]*workerInfo),
		roundRobinIndex: 0,
	}
}

func (s *service) RegisterWorker(workerID uint64, address string) error {
	s.workerPoolMutex.Lock()
	defer s.workerPoolMutex.Unlock()

	worker, ok := s.workerPool[workerID]
	if ok {
		worker.heartBeatMisses = 0
	} else {
		log.Println("Registering worker : ", workerID)

		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("Error in connecting to the worker at address : %s \nError : %v", address, err)
		}

		s.workerPool[workerID] = &workerInfo{
			address:       address,
			Connection:    conn,
			serviceClient: pb.NewWorkerServiceClient(conn),
		}

		s.workerPoolKeysMutex.Lock()
		defer s.workerPoolKeysMutex.Unlock()

		s.workerPoolKeys = make([]uint64, 0, len(s.workerPool))
		for k := range s.workerPool {
			s.workerPoolKeys = append(s.workerPoolKeys, k)
		}

		fmt.Println("Registered worker : ", workerID)

	}
	return nil
}

func (s *service) ManageWorkerPool(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.removeInactiveWorkers()
		case <-ctx.Done():
			return
		}
	}

}

func (s *service) ScanDatabase(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Performing Database Scan")
			go s.ExecuteAllScheduledTasks()
		case <-ctx.Done():
			log.Println("Shutting down database scanner.")
			return
		}
	}
}

func (s *service) ExecuteAllScheduledTasks() error {
	// 30 seconds for this transactions
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.repo.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
		/* 1. Get all tasks scheduled +30 seconds from now */
		tasks, err := s.repo.FetchScheduledTasks(ctx, tx)

		if err != nil {
			return err
		}

		/* 2. Send the task to the workers to process and update the task (picked_at) in the database */
		for _, task := range tasks {

			if err := s.submitTaskToWorker(task.ID, task.Command); err != nil {
				log.Printf("Failed to submit task %d: %v\n", task.ID, err)
				continue
			}
			// Dont come here unless task is sent successfully thats why continue
			if err := s.repo.MarkTaskAsPicked(ctx, tx, task.ID); err != nil {
				log.Printf("Failed to mark task %d as picked: %v\n", task.ID, err)
				continue
			}
		}

		return nil

	})

	if err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
		return err
	}

	return nil
}

func (s *service) UpdateTaskStatus(ctx context.Context, taskID uint64, status pb.TaskStatus, startedAt int64, completedAt int64, failedAt int64) error {

	var timeStamp time.Time
	var column string

	switch status {
	case pb.TaskStatus_STARTED:
		timeStamp = time.Unix(startedAt, 0)
		column = "started_at"
	case pb.TaskStatus_COMPLETED:
		timeStamp = time.Unix(completedAt, 0)
		column = "completed_at"
	case pb.TaskStatus_FAILED:
		timeStamp = time.Unix(failedAt, 0)
		column = "failed_at"
	default:
		log.Println("Invalid Status in UpdateStatusRequest")
		return errors.ErrUnsupported
	}

	err := s.repo.UpdateTaskStatus(ctx, taskID, timeStamp, column)

	return err
}

/* PRIVATE FUNCTIONS */

func (s *service) submitTaskToWorker(taskID uint64, data string) error {
	worker := s.getNextWorker()

	if worker == nil {
		return errors.New("no workers available")
	}

	// Which context?
	res, err := worker.serviceClient.ProcessTask(context.Background(), &pb.NewTaskRequest{
		TaskId: taskID,
		Data:   data,
	})
	log.Println("GRPC Response : ", res)
	return err
}

func (s *service) removeInactiveWorkers() {
	// Max 4
	s.workerPoolMutex.Lock()
	defer s.workerPoolMutex.Unlock()

	for workerID, worker := range s.workerPool {
		if worker.heartBeatMisses > 4 {
			log.Println("Removing inactive worker : ", workerID)
			worker.Connection.Close()
			delete(s.workerPool, workerID)

			s.workerPoolKeysMutex.Lock()

			// Rebuild keys slice for round robin
			s.workerPoolKeys = make([]uint64, 0, len(s.workerPool))
			for k := range s.workerPool {
				s.workerPoolKeys = append(s.workerPoolKeys, k)
			}

			s.workerPoolKeysMutex.Unlock()
		} else {

			// Every 4 Seconds the heartbeatmisses increases
			// If a worker doesnt notify the coordinator at max for 4 + 4 + 4 + 4 + 4 = 20 seconds then its declared dead
			worker.heartBeatMisses++
		}
	}
}

func (s *service) getNextWorker() *workerInfo {
	s.workerPoolMutex.Lock()
	defer s.workerPoolMutex.Unlock()

	workerCount := len(s.workerPoolKeys)
	if workerCount == 0 {
		return nil
	}

	worker := s.workerPool[s.workerPoolKeys[s.roundRobinIndex%uint32(workerCount)]]
	s.roundRobinIndex++
	return worker

}
