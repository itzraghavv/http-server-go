package request

import (
	"errors"
	"io"
	"strings"
)

type state int

const (
	initialized state = iota
	done
)

type Request struct {
	RequestLine RequestLine
	State       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var bufferSize int = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		State: initialized,
	}

	buff := make([]byte, bufferSize)
	readToIndex := 0

	for req.State != done {
		if readToIndex == len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		r, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		readToIndex += r

		bytesConsumed, err := req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		if bytesConsumed == 0 {
			continue
		}

		copy(buff, buff[bytesConsumed:readToIndex])
		readToIndex -= bytesConsumed
	}

	if req.State != done {
		return nil, errors.New("incomplete request")
	}

	return req, nil

}

func parseRequestLine(data string) (RequestLine, int, error) {

	idx := strings.Index(data, "\r\n")
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	requestLine := data[:idx]
	parts := strings.Split(requestLine, " ")

	bytesConsumed := idx + 2

	if len(parts) != 3 {
		return RequestLine{}, bytesConsumed, errors.New("invalid request parts")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return RequestLine{}, bytesConsumed, errors.New("invalid method")
		}
	}

	target := parts[1]

	version := parts[2]
	rawVer := strings.TrimPrefix(version, "HTTP/")
	if rawVer != "1.1" {
		return RequestLine{}, bytesConsumed, errors.New("unsupported HTTP version")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   rawVer,
	}, bytesConsumed, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == initialized {
		rl, noOfBytes, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if noOfBytes == 0 {
			return 0, nil
		}

		r.RequestLine = rl
		r.State = done

		return noOfBytes, nil
	}

	if r.State == done {
		return 0, errors.New("error: trying to read data in a done state")
	}

	return 0, errors.New("error: unknown state")
}
