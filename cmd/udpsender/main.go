package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("failed to resolve upd addr error: %s\n", err.Error())
	}

	udpConn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatalf("failed connection, error: %s\n", err.Error())
	}
	defer udpConn.Close()

	fmt.Printf("Sending to %s \n", "42069")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("> ")

		str, readErr := reader.ReadString('\n')
		if readErr != nil {
			log.Fatalf("failed to read, error: %s\n", readErr.Error())
		}

		_, writeErr := udpConn.Write([]byte(str))
		if writeErr != nil {
			log.Fatalf("failed to write, error: %s\n", writeErr.Error())
		}

		fmt.Printf("Message sent: %s", str)
	}
}
