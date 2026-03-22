package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gaganpreet-S1ngh/dts-ts/internal/platform/config"
	"github.com/Gaganpreet-S1ngh/dts-ts/internal/platform/database"
	server "github.com/Gaganpreet-S1ngh/dts-ts/internal/platform/httpserver"
	"github.com/Gaganpreet-S1ngh/dts-ts/internal/tasks"
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

	err = database.RegisterModels(dbClient, rootCtx, (*tasks.Task)(nil))

	// HTTP , GRPC
	httpServer := server.NewHTTPServer(cfg.ServerPort)

	// Dependency Injections
	repo := tasks.NewRepository(dbClient)
	svc := tasks.NewService(repo)
	handler := tasks.NewHandler(svc)
	routes := tasks.NewRoutes(handler)

	/* INITIALIZE MAIN APP */

	routes.SetupPublicRoutes()

	/* Start HTTP Server on a go routine */
	go func() {
		fmt.Println("Server Listening on port 8080")
		err := httpServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Error: %v", err)
		}

	}()

	/* SHUTDOWN MECHANISM */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal recieved")

	// Stop Server
	shutDownHTTP(httpServer)

	// Stop Workers
	rootCancel()

	// Close Services
	if err := dbClient.Close(); err != nil {
		log.Printf("DB Close Error: %v", err)
	}

}

func shutDownHTTP(s *http.Server) {
	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := s.Shutdown(ctx)

		if err != nil {
			log.Printf("Shutdown Error: %v", err)
			return
		}
	}
	log.Println("Server gracefully shutdown")
}
