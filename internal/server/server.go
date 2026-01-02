package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/itzraghavv/httpWebServer/internal/request"
	"github.com/itzraghavv/httpWebServer/internal/response"
)

type Server struct {
	Listener net.Listener
	handler  Handler
	isClosed atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(res *response.Writer, req *request.Request)

func (h *HandlerError) writeErr(w io.Writer) {
	if h == nil {
		return
	}

	body := []byte(h.Message)

	response.WriteStatusLine(w, h.StatusCode)
	headers := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(w, headers)
	w.Write(body)
}

func Serve(port int, handler Handler) (*Server, error) {
	listner, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		handler:  handler,
		Listener: listner,
	}

	go server.Listen()

	return server, nil
}

func (s *Server) Listen() {

	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue

		}
		go s.Handle(conn)
	}
}

var handlerErr Handler

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		hErr.writeErr(conn)
		return
	}

	buff := bytes.NewBuffer([]byte{})
	res := response.NewWriter(conn)
	s.handler(res, req)

	b := buff.Bytes()
	response.WriteStatusLine(conn, response.OK)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)

	conn.Write(b)

}

func (s *Server) Close() {
	s.Listener.Close()
	s.isClosed.Store(true)
}
