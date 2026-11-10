package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harsh0r/student-api/internal/config"
	"github.com/harsh0r/student-api/internal/http/handlers/student"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
	})
	router.HandleFunc("POST /api/student", student.New())
	// setup server
	server := http.Server{Addr: cfg.Address, Handler: router}
	fmt.Printf("server started at %s.", cfg.Address)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server.")
		}
	}()
	<-done
	slog.Info("shutting down server.")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("Error: ", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
