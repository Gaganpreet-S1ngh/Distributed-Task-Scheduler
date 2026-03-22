package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gaganpreet-S1ngh/dts-proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NewTaskInfo struct {
	taskID uint64
	data   string
	status pb.TaskStatus
}

type Service interface {
	ConnectToCoordinator(coordinatorServerAddress string) error
	StartWorkerPool(ctx context.Context, numWorkers int)
	AddToTaskQueue(task *NewTaskInfo)
	PeriodicHeartbeat(ctx context.Context, workerAddress string)
}

type service struct {
	ID            uint64
	serviceClient pb.CoordinatorServiceClient
	mu            sync.Mutex
	Tasks         map[uint64]*NewTaskInfo
	taskQueue     chan *NewTaskInfo
	wg            sync.WaitGroup
}

func NewService() Service {
	return &service{
		//use uuid
		ID:        1,
		Tasks:     make(map[uint64]*NewTaskInfo),
		taskQueue: make(chan *NewTaskInfo, 100),
	}
}

func (s *service) PeriodicHeartbeat(ctx context.Context, workerAddress string) {
	s.wg.Add(1)
	defer s.wg.Done()

	// Send heartbeat every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.sendHeartBeat(workerAddress); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)

			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *service) ConnectToCoordinator(coordinatorServerAddress string) error {
	log.Println("Connecting to Coordinator... ")
	conn, err := grpc.NewClient(coordinatorServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("Error connecting to coordinator : %v", err)
	}

	s.serviceClient = pb.NewCoordinatorServiceClient(conn)
	log.Println("Connected to coordinator!")
	return nil

}

func (s *service) StartWorkerPool(ctx context.Context, numWorkers int) {
	for range numWorkers {
		s.wg.Add(1)
		go s.worker(ctx)
	}
}

func (s *service) AddToTaskQueue(task *NewTaskInfo) {
	// Check for duplicate task?

	s.mu.Lock()
	s.Tasks[task.taskID] = task
	s.mu.Unlock()

	s.taskQueue <- task

}

/* PRIVATE FUNCTIONS */

func (s *service) worker(ctx context.Context) {
	defer s.wg.Done()

	// keep looking for new tasks from the channel
	for {
		select {
		case task := <-s.taskQueue:
			//Recieved task update tell the coordinator the task has started (A Call with no cpu time required hence handled by a go routine)
			go s.updateTaskStatus(task, pb.TaskStatus_STARTED)
			// Main processing done by the CPU
			s.processTask(task)
			// After completing the task send the status of completed
			go s.updateTaskStatus(task, pb.TaskStatus_COMPLETED)

		case <-ctx.Done():
			return
		}
	}

}

func (s *service) updateTaskStatus(task *NewTaskInfo, status pb.TaskStatus) {
	// Update local map
	s.mu.Lock()
	task.status = status
	s.mu.Unlock()

	// Send both started at and completed at as now because it is handled in our coordinator service according to the status
	req := &pb.UpdateTaskStatusRequest{
		TaskId:      task.taskID,
		Status:      status,
		StartedAt:   time.Now().Unix(),
		CompletedAt: time.Now().Unix(),
	}

	// Which context to use in grpcClient calls?

	s.serviceClient.UpdateTaskStatus(context.Background(), req)
}

func (s *service) processTask(task *NewTaskInfo) {

	// Simulating processing of task
	log.Printf("Processing task: %+v", task)
	time.Sleep(2 * time.Second)
	task.data = "This task was completed by the worker"
	log.Printf("Completed task: %+v", task)
}

func (s *service) sendHeartBeat(workerAddress string) error {
	res, err := s.serviceClient.SendHeartbeat(context.Background(), &pb.HeartbeatRequest{
		WorkerId: s.ID,
		Address:  workerAddress,
	})

	log.Println("GRPC Response : ", res)
	return err
}
