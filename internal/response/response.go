package response

import (
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
