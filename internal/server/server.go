package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/itzraghavv/httpWebServer/internal/response"
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

func (s *Server) Handle(conn net.Conn) error {
	defer conn.Close()

	body := []byte("")

	headers := response.GetDefaultHeaders(len(body))

	err := response.WriteStatusLine(conn, response.OK)
	if err != nil {
		return err
	}

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return err
	}

	conn.Write(body)
	return nil
}

func (s *Server) Close() {
	s.Listener.Close()
	s.isClosed.Store(true)
}
