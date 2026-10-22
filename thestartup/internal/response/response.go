package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	StatusOk            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

type writerStatus int

const (
	writerStatusLine    writerStatus = iota
	writerStatusHeaders writerStatus = iota
	writerStatusBody    writerStatus = iota
)

type Writer struct {
	io.Writer
	status writerStatus
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: w,
		status: writerStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.status != writerStatusLine {
		return errors.New("Request out of order")
	}
	reason := ""
	switch statusCode {
	case StatusOk:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalError:
		reason = "Internal Server Error"
	default:
		reason = "OK"
	}
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, reason)
	w.status = writerStatusHeaders
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = fmt.Sprintf("%d", contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.status != writerStatusHeaders {
		return errors.New("Request out of order")
	}
	if len(headers) == 0 {
		return errors.New("Headers section empty")
	}
	for k, v := range headers {
		fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w, "\r\n")
	w.status = writerStatusBody
	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.status != writerStatusBody {
		return 0, errors.New("Request out of order")
	}
	return w.Writer.Write(body)
}
func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	length := len(p)

	n, err := fmt.Fprintf(w.Writer, "%x\r\n", length)
	if err != nil {
		return n, err
	}
	return w.Writer.Write(append(p, []byte("\r\n")...))
}
func (w *Writer) WriteChunkedBodyDone(hasTrailer bool) (int, error) {
	if !hasTrailer {
		return fmt.Fprintf(w.Writer, "0\r\n\r\n")
	} else {
		return fmt.Fprintf(w.Writer, "0\r\n")
	}
}
func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		_, err := fmt.Fprintf(w.Writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w.Writer, "\r\n")
	return err
}
