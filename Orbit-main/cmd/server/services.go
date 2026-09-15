package main

import (
	"Orbit/internal/db"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GraceFulShutDown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// once it will get signal only then from here will start executing..

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	ExecuteShutDown(ctx, server)
}

func ExecuteShutDown(ctx context.Context, server *http.Server) {
	stopIncomingRequest(server)
	stopDatabase()
	stopRedis()
}

func stopIncomingRequest(server *http.Server) {
	tempContext, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(tempContext); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func stopDatabase() {
	client := db.GetClient()
	if client != nil {
		tempContext, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := client.Disconnect(tempContext); err != nil {
			log.Printf("Database Disconnect Error")
		}
	}
}

func stopRedis() {
	client := db.GetRedisClient()
	if client != nil {
		if err := client.Close(); err != nil {
			log.Printf("Redis Disconnect Error")
		}
	}
}
