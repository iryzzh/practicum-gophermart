package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	handler       http.Handler
	serverAddress string
}

func New(serverAddress string, handler http.Handler) *Server {
	return &Server{
		serverAddress: serverAddress,
		handler:       handler,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	listener, err := s.getListener()
	if err != nil {
		return err
	}
	defer listener.Close()

	srv := &http.Server{
		Handler:        s.handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Starting HTTP server on", listener.Addr())

	serveError := make(chan error, 1)
	go func() {
		select {
		case serveError <- srv.Serve(listener):
		case <-ctx.Done():
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down HTTP server")
	case err := <-serveError:
		log.Println("HTTP server error:", err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeout); err == timeout.Err() {
		srv.Close()
	}

	return nil
}

func (s *Server) getListener() (net.Listener, error) {
	l, err := net.Listen("tcp", s.serverAddress)
	if err != nil {
		return nil, err
	}

	return l, nil
}
