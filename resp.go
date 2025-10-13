package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseRESP(buf *bufio.Reader) (res []string, err error) {
	parsed := []string{}

	line, err := buf.ReadString('\n')

	if err != nil {
		if err == io.EOF {
			return parsed, fmt.Errorf("unexpected end of input")
		}
		return
	}

	line = strings.TrimSuffix(line, "\r\n")

	switch line[0] {
	case '$':
		expectedLength, err := strconv.Atoi(line[1:])
		if err != nil || expectedLength < 0 {
			return nil, fmt.Errorf("invalid bulk string length")
		}

		nextLine, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("unexpected end of input")
			}
			return nil, err
		}

		nextLine = strings.TrimSuffix(nextLine, "\r\n")

		if len(nextLine) != expectedLength {
			return nil, fmt.Errorf("bulk string length mismatch")
		}

		parsed = append(parsed, nextLine)
	case '*':
		expectedLength, err := strconv.Atoi(line[1:])
		if err != nil || expectedLength < 0 {
			return nil, fmt.Errorf("invalid array length")
		}

		for i := 0; i < expectedLength; i++ {
			subRes, err := parseRESP(buf)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, subRes...)
		}
	}

	return parsed, nil
}
