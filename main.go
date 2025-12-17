package main

import (
	"context"
	"db_labs/ioc"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	ShutdownPeriod     = 15 * time.Second
	HardShutdownPeriod = 3 * time.Second
)

func main() {
	db := ioc.UsePgConnection()
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := godotenv.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading .env file: %v", err))
	}

	ongoingCtx, stopOngoing := context.WithCancel(context.Background())

	server := ioc.UseHttpServer(ongoingCtx)

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownPeriod)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	stopOngoing()
	if err != nil {
		slog.Error("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(HardShutdownPeriod)
	}
}
