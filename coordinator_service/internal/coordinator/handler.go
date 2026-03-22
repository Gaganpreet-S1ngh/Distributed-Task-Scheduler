package coordinator

import (
	"context"
	"log"

	"github.com/Gaganpreet-S1ngh/dts-proto/pb"
	"google.golang.org/grpc"
)

type Handler struct {
	pb.UnimplementedCoordinatorServiceServer
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

/* Register Handler to Handler GRPC Calls */
func (h *Handler) Register(srv *grpc.Server) {
	pb.RegisterCoordinatorServiceServer(srv, h)
}

/* GRPC */

func (h *Handler) SubmitTask(ctx context.Context, input *pb.FinishedTaskRequest) (*pb.FinishedTaskResponse, error) {
	// Simulate Recieving of completed tasks -> Parsing and storing in separate tables

	return &pb.FinishedTaskResponse{Message: "Task Successfully Submitted"}, nil
}

func (h *Handler) UpdateTaskStatus(ctx context.Context, input *pb.UpdateTaskStatusRequest) (*pb.UpdateTaskStatusResponse, error) {
	err := h.svc.UpdateTaskStatus(ctx, input.GetTaskId(), input.GetStatus(), input.GetStartedAt(), input.GetCompletedAt(), input.GetFailedAt())

	if err != nil {
		log.Println("Error in updating task : Error ", err)
		return nil, err
	}

	return &pb.UpdateTaskStatusResponse{Success: true}, nil
}

func (h *Handler) SendHeartbeat(ctx context.Context, input *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	err := h.svc.RegisterWorker(uint64(input.GetWorkerId()), input.GetAddress())
	if err != nil {
		return nil, err
	}

	return &pb.HeartbeatResponse{Acknowledged: true}, nil
}
