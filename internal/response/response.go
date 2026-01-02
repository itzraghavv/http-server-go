package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/itzraghavv/httpWebServer/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

type writerState string

func writeStateInit() writerState {
	return "writer initialized"
}

func writerStateStatusLineDone() writerState {
	return "StatusLine Written"
}

func writerStateHeadersDone() writerState {
	return "Headers Written"
}

func writerStateDone() writerState {
	return "done"
}

type Writer struct {
	conn  io.Writer
	state writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writeStateInit() {
		return fmt.Errorf("status line writer started")
	}
	// buf := []byte{}
	// _, err := w.conn.Write(buf)
	// if err == nil {
	// 	w.state = writerStateStatusLineDone()
	// }
	// return err
	var line string
	switch statusCode {
	case OK:
		line = "HTTP/1.1 200 OK\r\n"
	case BadRequest:
		line = "HTTP/1.1 400 Bad Request\r\n"
	case InternalServerError:
		line = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		line = "HTTP/1.1 " + strconv.Itoa(int(statusCode)) + "\r\n"
	}

	if _, err := w.conn.Write([]byte(line)); err != nil {
		return err
	}

	w.state = writerStateStatusLineDone()
	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != writerStateStatusLineDone() {
		return fmt.Errorf("cannot write headers before status line")
	}
	for key, value := range h {
		line := key + ": " + value + "\r\n"
		if _, err := w.conn.Write([]byte(line)); err != nil {
			return err
		}
	}

	// blank line between headers and body
	if _, err := w.conn.Write([]byte("\r\n")); err != nil {
		return err
	}

	w.state = writerStateHeadersDone()
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateHeadersDone() {
		return 0, fmt.Errorf("trying to WriteBody() not in HeadersWritten state")
	}
	n, err := w.conn.Write(p)
	if err == nil {
		w.state = writerStateDone()
		w.state = writerStateDone()
	}
	return n, err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:

		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}

	case BadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}

	case InternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}

	default:
		_, err := w.Write([]byte("HTTP/1.1" + strconv.Itoa(int(statusCode)) + "\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	resHeader := make(headers.Headers)

	resHeader["Content-Length"] = strconv.Itoa(contentLen)
	resHeader["Connection"] = "close"
	resHeader["Content-Type"] = "text/plain"

	return resHeader
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		line := key + ": " + value + "\r\n"

		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	return nil
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		state: writeStateInit(),
		conn:  conn,
	}
}
