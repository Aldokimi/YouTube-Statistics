package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

// "github.com/joho/godotenv"


func newRouter() *httprouter.Router {
	mux := httprouter.New()

	// load the env variables
	// err := godotenv.Load(); if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	
	ytApiKey := os.Getenv("ytAPIKey"); if  ytApiKey == "" {
		log.Fatal("youtube API key not provided")
	}

	mux.GET("/sum", getChannelsStatus(ytApiKey))
	return mux
}

func main() {
	srv := &http.Server{
		Addr: ":8000",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Println("shutdown complete")

		close(idleConnsClosed)
	}()

	log.Printf("Starting server on port 8000")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Server Stopped!")
}