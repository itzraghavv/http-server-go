package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	data := string(r)

	lines := strings.Split(data, "\r\n")

	requestLine := lines[0]

	rl, err := parseRequestLine(requestLine)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: rl,
	}, nil

}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, errors.New("invalid request parts")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return RequestLine{}, errors.New("invalid method")
		}
	}

	target := parts[1]

	version := parts[2]
	rawVer := strings.TrimPrefix(version, "HTTP/")
	if rawVer != "1.1" {
		return RequestLine{}, errors.New("unsupported HTTP version")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   rawVer,
	}, nil
}
