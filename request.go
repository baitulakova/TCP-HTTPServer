package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Request struct {
	Method   string
	URL      *url.URL
	Protocol string
	Headers  map[string][]string
	Body     []byte
}

func readLinebyLine(filepath string) ([]string, error) {
	lines := make([]string, 0)

	file, err := os.Open(filepath)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// reads request lines and converts to Request struct
func convertToRequest(filepath string) (*Request, error) {
	lines, err := readLinebyLine(filepath)
	if err != nil {
		return nil, err
	}
	req := &Request{}

	req.Headers = make(map[string][]string)
	if len(lines) == 0 {
		return nil, errors.New("empty request body")
	}

	firstLine := strings.Split(lines[0], " ") // method, URL, protocol version

	if len(firstLine) != 3 {
		return nil, errors.New("invalid first line of request body")
	}
	req.Method = firstLine[0]

	u, err := url.Parse(firstLine[1])
	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %v", err)
	}

	req.URL = u
	req.Protocol = firstLine[2]

	lines = lines[1:]
	if len(lines) == 0 {
		return nil, errors.New("invalid request body")
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			keyHeader := headerParts[0]
			valueHeader := headerParts[1]

			var values []string
			values = append(values, valueHeader)
			req.Headers[keyHeader] = values

			req.Body = []byte(headerParts[0])
		} else {
			req.Body = []byte(line)
		}
	}

	return req, nil
}
