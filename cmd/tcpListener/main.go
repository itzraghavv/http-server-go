package main

import (
	"fmt"
	"log"
	"net"

	glc "github.com/itzraghavv/httpWebServer/getLinesChannel"
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
		fmt.Println("Connection accepted", conn.RemoteAddr())

		l := glc.GetLinesChannel(conn)
		for line := range l {
			fmt.Println(line)
		}

		fmt.Println("connection closed", conn.RemoteAddr())
	}

}
