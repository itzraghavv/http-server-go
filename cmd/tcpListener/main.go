package main

import (
	"fmt"
	"log"
	"net"

	"github.com/itzraghavv/httpWebServer/internal/request"
)

func main() {

	tcpListner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}

	defer tcpListner.Close()

	for {
		conn, conErr := tcpListner.Accept()
		if conErr != nil {
			log.Fatalf("error: %s\n", conErr.Error())
		}
		log.Println("Connection accepted", conn.RemoteAddr())

		data, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Printf("Request Line:\n - Method: %v\n - Target: %v\n - Version: %v\n", data.RequestLine.Method, data.RequestLine.RequestTarget, data.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for k, value := range data.Headers {
			fmt.Printf("- %s: %s\n", k, value)
		}
		fmt.Printf("Body:\n - %v", string(data.Body))
		log.Println("connection closed", conn.RemoteAddr())
	}
}
