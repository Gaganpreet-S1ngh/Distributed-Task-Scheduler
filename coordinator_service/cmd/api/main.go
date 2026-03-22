package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gaganpreet-S1ngh/dts-cooridnator/internal/coordinator"
	"github.com/Gaganpreet-S1ngh/dts-cooridnator/internal/platform/config"
	"github.com/Gaganpreet-S1ngh/dts-cooridnator/internal/platform/database"
	"github.com/Gaganpreet-S1ngh/dts-cooridnator/internal/platform/grpcserver"
)

func main() {

	/* INITIALIZE ENV AND CONTEXTS */

	cfg := config.LoadConfig()
	rootCtx, rootCancel := context.WithCancel(context.Background())

	/* INITIALIZE SERVICES */

	// Database
	dbClient, err := database.ConnectToDB(cfg.DatabaseDSN, rootCtx)

	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// HTTP , GRPC
	grpcServer, err := grpcserver.NewGrpcServer(cfg.GrpcServerPort)

	if err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}

	// Dependency Injections

	repo := coordinator.NewRepository(dbClient)
	svc := coordinator.NewService(repo)
	handler := coordinator.NewHandler(svc)

	/* INITIALIZE MAIN APP */

	// Start Manage worker pool in a go routine

	go svc.ManageWorkerPool(rootCtx)

	// Start GRPC Server on a go routine
	handler.Register(grpcServer.Server)
	go func() {
		log.Printf("gRPC listening on :%s", cfg.GrpcServerPort)
		if err := grpcServer.Server.Serve(grpcServer.Listener); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	// Keep scanning for work in database in a go routine
	go svc.ScanDatabase(rootCtx)

	/* SHUTDOWN MECHANISM */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal recieved")

	// Stop Server
	grpcServer.Server.GracefulStop()

	// Stop Workers
	rootCancel()

	// Close Services
	if err := dbClient.Close(); err != nil {
		log.Printf("DB Close Error: %v", err)
	}

}
