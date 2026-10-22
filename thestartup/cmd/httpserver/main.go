package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func sendVideoResponse(w *response.Writer, _ *request.Request) {
	file, _ := os.ReadFile("./assets/vim.mp4")
	headers := response.GetDefaultHeaders(len(file))
	headers["content-type"] = "video/mp4"
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(headers)
	w.WriteBody(file)
}
func createResponse(c response.StatusCode) string {
	title := ""
	h1 := ""
	p := ""
	switch c {
	case response.StatusOk:
		title = "200 OK"
		h1 = "Success!"
		p = "Your request was an absolute banger."
	case response.StatusBadRequest:
		title = "400 Bad Request"
		h1 = "Bad Request"
		p = "Your request honestly kinda sucked."
	case response.StatusInternalError:
		title = "500 Internal Server Error"
		h1 = "Internal Server Error"
		p = "Okay, you know what? This one is on me."
	}

	return fmt.Sprintf(`<html><head><title>%s</title></head><body><h1>%s</h1><p>%s</p></body></html>`, title, h1, p)
}
func streamResponse(w *response.Writer, _ *request.Request, cuantity int64) {
	headers := response.GetDefaultHeaders(0)
	delete(headers, "content-length")
	headers["transfer-encoding"] = "chunked"
	headers["content-type"] = "text/plain"
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(headers)
	resp, _ := http.Get(fmt.Sprintf("https://httpbingo.org/stream/%d", cuantity))
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil || n == 0 {
			w.WriteChunkedBodyDone(false)
		}
		w.WriteChunkedBody(buffer[:n])
	}

}

func streamRangeResponse(w *response.Writer, _ *request.Request, cuantity int64) {
	h := response.GetDefaultHeaders(0)
	delete(h, "content-length")
	h["transfer-encoding"] = "chunked"
	h["content-type"] = "text/plain"
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)
	resp, _ := http.Get(fmt.Sprintf("https://httpbingo.org/range/%d", cuantity))
	fullBuffer := []byte{}
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			fullBuffer = append(fullBuffer, buffer[:n]...)
			w.WriteChunkedBody(buffer[:n])
		}
		if err != nil {
			break
		}
	}
	w.WriteChunkedBodyDone(true)
	trailerH := headers.NewHeaders()
	hash := sha256.Sum256(fullBuffer)
	trailerH["x-content-sha256"] = fmt.Sprintf("%x", hash)
	trailerH["x-content-length"] = fmt.Sprintf("%d", len(fullBuffer))
	w.WriteTrailers(trailerH)

}

func handler(w *response.Writer, r *request.Request) {
	if after, ok := strings.CutPrefix(r.RequestLine.RequestTarget, "/httpbin/stream/"); ok {
		if i, err := strconv.ParseInt(after, 10, 64); err == nil {
			streamResponse(w, r, i)
			return
		}
	}
	if after, ok := strings.CutPrefix(r.RequestLine.RequestTarget, "/httpbin/range/"); ok {
		if i, err := strconv.ParseInt(after, 10, 64); err == nil {
			streamRangeResponse(w, r, i)
			return
		}
	}
	body := ""
	var statusCode response.StatusCode
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		body = createResponse(response.StatusBadRequest)
		statusCode = response.StatusBadRequest
	case "/myproblem":
		body = createResponse(response.StatusInternalError)
		statusCode = response.StatusInternalError
	case "/video":
		sendVideoResponse(w, r)
		return
	default:
		body = createResponse(response.StatusOk)
		statusCode = response.StatusOk
	}
	w.WriteStatusLine(statusCode)
	h := response.GetDefaultHeaders(len(body))
	h["content-type"] = "text/html"
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))

}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
