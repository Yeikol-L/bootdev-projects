package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	Message string
	Code    response.StatusCode
}

func (err *HandlerError) writeToReq(w *response.Writer) {
	headers := response.GetDefaultHeaders(len(err.Message))
	w.WriteStatusLine(response.StatusCode(err.Code))
	w.WriteHeaders(headers)
	w.Write([]byte(err.Message))
}

type Handler func(w *response.Writer, req *request.Request)

// type Server struct - Contains the state of the server
//
// func Serve(port int) (*Server, error) - Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
//
// func (s *Server) Close() error - Closes the listener and the server
//
// func (s *Server) listen() - Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine. I used an atomic.Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.
//
// func (s *Server) handle(conn net.Conn) - Handles a single connection by writing the following response and then closing the connection:
type Server struct {
	port     int
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := Server{port: port, listener: socket, handler: handler}
	go s.listen()
	return &s, nil
}

func (s *Server) Close() {
	if s.listener != nil {
		s.closed.Store(true)
		s.listener.Close()
	}
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		h := HandlerError{Message: err.Error(), Code: response.StatusBadRequest}
		h.writeToReq(response.NewWriter(conn))
		return
	}
	writer := response.NewWriter(conn)
	s.handler(writer, req)
}
