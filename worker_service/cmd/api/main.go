package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gaganpreet-S1ngh/dts-worker/internal/platform/config"
	"github.com/Gaganpreet-S1ngh/dts-worker/internal/platform/grpcserver"
	"github.com/Gaganpreet-S1ngh/dts-worker/internal/worker"
)

func main() {

	/* INITIALIZE ENV AND CONTEXTS */

	cfg := config.LoadConfig()
	rootCtx, rootCancel := context.WithCancel(context.Background())

	/* INITIALIZE SERVICES */

	// HTTP , GRPC
	grpcServer, err := grpcserver.NewGrpcServer(cfg.GrpcServerPort)

	if err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}

	// Dependency Injections
	svc := worker.NewService()
	handler := worker.NewHandler(svc)

	/* INITIALIZE MAIN APP */

	// Start the worker pool
	svc.StartWorkerPool(rootCtx, 4)

	// Connect to coordinator
	svc.ConnectToCoordinator(cfg.GrpcServerAddress)

	// Start sending heartbeats to the coordinator

	go svc.PeriodicHeartbeat(rootCtx, grpcServer.Listener.Addr().String())

	// Start gRPC Server in a goroutine
	handler.Register(grpcServer.Server)
	go func() {
		actualPort := grpcServer.Listener.Addr().(*net.TCPAddr).Port
		log.Printf("gRPC listening on :%d", actualPort)
		if err := grpcServer.Server.Serve(grpcServer.Listener); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	/* SHUTDOWN MECHANISM */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal recieved")

	// Stop Server
	grpcServer.Server.GracefulStop()

	// Stop Workers
	rootCancel()

}
