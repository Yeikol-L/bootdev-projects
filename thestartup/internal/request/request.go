package request

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

const bufferSize = 8

const crlf = "\r\n"

type requestState int

const (
	requestStateInitialized    requestState = iota
	requestStateParsingReqLine requestState = iota
	requestStateParsingHeaders requestState = iota
	requestStateParsingBody    requestState = iota
	requestStateDone           requestState = iota
)

type Request struct {
	RequestLine  RequestLine
	requestState requestState
	headers.Headers
	Body []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) done() bool {
	return r.requestState == requestStateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		requestState: requestStateInitialized,
		Headers:      headers.NewHeaders(),
	}
	for !req.done() {
		if readToIndex == cap(buffer) {
			newBufer := make([]byte, cap(buffer)*2)
			copy(newBufer, buffer)
			buffer = newBufer
		}
		numBytesRead, buffErr := reader.Read(buffer[readToIndex:])
		if buffErr != nil {
			if errors.Is(buffErr, io.EOF) {
				if !req.done() {
					return nil, errors.New("incomplete request")
				}
				req.requestState = requestStateDone
				break
			}
			return nil, buffErr
		}
		readToIndex += numBytesRead
		numBytesRead, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}
		if numBytesRead > 0 {
			newBuffer := make([]byte, len(buffer))
			copy(newBuffer, buffer[numBytesRead:])
			buffer = newBuffer
			readToIndex -= numBytesRead
		}
	}
	return req, nil
}

func (r *Request) parseRequestLine(content []byte) (int, error) {
	contentString, _, found := strings.Cut(string(content), crlf)
	if !found {
		return 0, nil
	}

	parts := strings.Split(contentString, " ")
	if len(parts) != 3 {
		return 0, errors.New("Malformed request line")
	}
	r.RequestLine.Method = parts[0]
	r.RequestLine.RequestTarget = parts[1]
	version, found := strings.CutPrefix(parts[2], "HTTP/")
	if !found {
		return 0, errors.New("Malformed request line")
	}
	r.RequestLine.HttpVersion = version
	r.requestState = requestStateParsingHeaders
	return len(contentString) + 2, nil
}

func (r *Request) parseRequestHeader(content []byte) (int, error) {
	parsed, done, err := r.Headers.Parse(content)
	if err != nil {
		return 0, err
	}
	if done {
		if r.Headers.Get("Content-Length") == "" {
			r.requestState = requestStateDone
		} else {
			r.requestState = requestStateParsingBody
		}
	}
	return parsed, nil
}
func (r *Request) parseRequestBody(body []byte, contentLength int64) (int, error) {
	remainingLength := contentLength - int64(len(r.Body))
	if remainingLength < int64(len(body)) {
		return 0, fmt.Errorf("Content length mismatch %d - %d", len(body), remainingLength)
	}
	r.Body = append(r.Body, body...)
	if int64(len(r.Body)) == contentLength {
		r.requestState = requestStateDone
	}
	return len(body), nil
}

func (r *Request) parse(data []byte) (int, error) {
	total := 0
	contentLength := int64(0)
	if r.requestState == requestStateInitialized {
		r.requestState = requestStateParsingReqLine
	}
	for !r.done() {
		if r.requestState == requestStateParsingBody {
			v, err := strconv.ParseInt(r.Get("Content-Length"), 10, 64)
			if err == nil {
				contentLength = v
			}
		}
		n, err := r.parseOne(data[total:], contentLength)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		total += n
	}
	return total, nil
}

func (r *Request) parseOne(data []byte, contentLength int64) (int, error) {
	if r.requestState == requestStateParsingReqLine {
		parsed, err := r.parseRequestLine(data)
		if err != nil {
			r.requestState = requestStateDone
		}
		return parsed, err
	}
	if r.requestState == requestStateParsingHeaders {
		parsed, err := r.parseRequestHeader(data)
		if err != nil {
			r.requestState = requestStateDone
		}
		return parsed, err
	}
	if r.requestState == requestStateParsingBody {
		parsed, err := r.parseRequestBody(data, contentLength)
		if err != nil {
			r.requestState = requestStateDone
		}
		return parsed, err
	}
	return 0, nil
}
