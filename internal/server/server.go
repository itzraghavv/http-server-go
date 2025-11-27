package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listner, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		Listener: listner,
	}

	go server.Listen()

	return server, nil
}

func (s *Server) Listen() error {

	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			if s.isClosed.Load() {
				break
			} else {
				log.Println(err)
				continue
			}

		}
		go s.Handle(conn)
	}
	return nil
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"))
}

func (s *Server) Close() {
	s.isClosed.Store(true)
	s.Listener.Close()
}
