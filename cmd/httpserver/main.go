package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itzraghavv/httpWebServer/internal/request"
	"github.com/itzraghavv/httpWebServer/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {

	if req == nil {
		log.Println("req is nil")
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "nil request\n",
		}
	}

	target := req.RequestLine.RequestTarget

	if target == "/myproblem" {
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	}
	return &server.HandlerError{
		Message: "All good",
	}
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
