package server

import "time"

// Option configures server by overriding a default setting.
type Option func(*Server)

// WithTimeout sets the maximum duration for requests.
func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}
