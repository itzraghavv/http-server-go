package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

var tcharRegex = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+$`)

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	line := string(data[:idx])
	line = strings.TrimLeft(line, " \t")
	if line == "" {
		return 0, false, errors.New("empty header line")
	}

	cIdx := strings.Index(line, ":")
	if cIdx == -1 {
		return 0, false, errors.New("missing colon in header")
	}

	fieldName := line[:cIdx]
	fieldValue := strings.TrimSpace(line[cIdx+1:])

	if strings.ContainsAny(fieldName, " \t") {
		return 0, false, errors.New("invalid whitespace in header name")
	}

	if !tcharRegex.MatchString(fieldName) {
		return 0, false, errors.New("contains invalid characters")
	}

	fieldName = strings.ToLower(fieldName)

	existing, ok := h[fieldName]
	if ok {
		h[fieldName] = existing + ", " + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	return idx + 2, false, nil
}

func (h Headers) Get(key string) string {
	value := h[strings.ToLower(key)]

	return value
}
