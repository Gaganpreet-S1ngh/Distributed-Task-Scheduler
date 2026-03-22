package worker

import (
	"context"
	"log"

	"github.com/Gaganpreet-S1ngh/dts-proto/pb"
	"google.golang.org/grpc"
)

type Handler struct {
	pb.UnimplementedWorkerServiceServer
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

/* Register Handler to Handler GRPC Calls */
func (h *Handler) Register(srv *grpc.Server) {
	pb.RegisterWorkerServiceServer(srv, h)
}

/* GRPC */

func (h *Handler) ProcessTask(ctx context.Context, input *pb.NewTaskRequest) (*pb.NewTaskResponse, error) {
	log.Printf("Received task: %+v", input)
	h.svc.AddToTaskQueue(&NewTaskInfo{
		taskID: input.GetTaskId(),
		data:   input.GetData(),
	})

	return &pb.NewTaskResponse{
		TaskId:  input.GetTaskId(),
		Message: "Task recieved successfully for processing",
		Success: true,
	}, nil
}
