package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// defaultTimeout is the time to wait before timing out a connection.
	defaultTimeout = time.Second * 5
)

// Server accepts and handles TCP connections.
type Server struct {
	listener net.Listener
	logger   *slog.Logger
	timeout  time.Duration
	wg       sync.WaitGroup
}

// New initializes a server with optional configuration.
func New(logger *slog.Logger, port uint, opts ...Option) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("start tcp listener: %w", err)
	}

	s := &Server{
		listener: ln,
		logger:   logger,
		timeout:  defaultTimeout,
	}

	// apply optional configuration
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// ListenAndServe starts a server and blocks until the context is cancelled or
// the server is stopped. The server is gracefully stopped.
//
// Once it has been stopped it is NOT safe for reuse.
func (s *Server) ListenAndServe(ctx context.Context) {
	s.logger.Info("start accepting tcp conns", "addr", s.Addr())

	go func() {
		<-ctx.Done()
		s.Shutdown()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Info("waiting for clients to disconnect", "addr", s.Addr())
				s.wg.Wait()
				s.logger.Info("shutting down server", "addr", s.Addr())
				return
			}
			s.logger.Error("accept conn", "error", err)
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConn(conn)
		}()
	}
}

// Shutdown stops the server from accepting new connections and causes any
// blocked operations to return.
func (s *Server) Shutdown() {
	s.logger.Info("shutdown signal received", "addr", s.Addr())
	s.logger.Info("stop accepting tcp conns", "addr", s.Addr())
	s.listener.Close()
}

// Addr returns the listener's network address.
func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

// handleConn accepts and handles a TCP connection. If any errors are
// encountered, the connection is closed.
func (s *Server) handleConn(conn net.Conn) {
	clientID, _ := uuid.NewRandom()
	s.logger.Info("accepted conn", "clientID", clientID, "localAddr", conn.LocalAddr(), "remoteAddr", conn.RemoteAddr())

	defer func() {
		_ = conn.Close()
		s.logger.Info("closed conn", "clientID", clientID)
	}()

	_ = conn.SetDeadline(time.Now().Add(s.timeout))

	if _, err := io.Copy(conn, conn); err != nil {
		s.logger.Error("responded", "clientID", clientID, "error", err)
		return
	}
}
