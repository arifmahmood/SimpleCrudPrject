package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple-crud-project/api"
	"simple-crud-project/repo"
	"syscall"
	"time"
)

// Server ...
type Server struct {
	urlRepo   repo.Url
}

// NewServer ...
func NewServer(url repo.Url) *Server {
	return &Server{
		urlRepo:   url,
	}
}

// Serve ...
func (s *Server) Serve() {

	portStr := viper.GetString("PORT")

	r := chi.NewMux()
	r.Mount("/api/v1", api.NewRouter(s.urlRepo))

	server := &http.Server{
		ReadTimeout:  viper.GetDuration("READ_TIMEOUT") * time.Second,
		WriteTimeout: viper.GetDuration("WRITE_TIMEOUT") * time.Second,
		IdleTimeout:  viper.GetDuration("IDLE_TIMEOUT") * time.Second,
		Addr:         ":" + portStr,
		Handler:      r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Println("Server Listening on :" + portStr)
		log.Fatal(server.ListenAndServe())
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Server shut down gracefully")
}


