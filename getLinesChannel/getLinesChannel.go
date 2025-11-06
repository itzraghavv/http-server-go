package getLineschannel

import (
	"io"
	"strings"
)

func GetLinesChannel(f io.ReadCloser) <-chan string {

	strCh := make(chan string)

	go func() {
		defer close(strCh)

		byteArr := make([]byte, 8)

		currLine := ""

		for {
			n1, err := f.Read(byteArr)
			if err != nil && err != io.EOF {
				return
			}

			part := string(byteArr[:n1])
			currLine = currLine + part
			parts := strings.Split(currLine, "\n")

			for _, p := range parts[:len(parts)-1] {
				strCh <- p
			}

			currLine = parts[len(parts)-1]

			if err == io.EOF {
				if currLine != "" {
					strCh <- currLine // send leftover
				}
				return
			}
		}
	}()

	return strCh
}
