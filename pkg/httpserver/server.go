package httpserver

import (
	"context"
	"net/http"
	"time"
)

// NB (alkurbatov): Set reasonable timeouts, see:
// https://habr.com/ru/company/ispring/blog/560032/
const (
	_defaultReadTimeout       = 5 * time.Second
	_defaultWriteTimeout      = 10 * time.Second
	_defaultIdleTimeout       = 120 * time.Second
	_defaultReadHeaderTimeout = 5 * time.Second
)

type Server struct {
	server *http.Server
	notify chan error
}

func New(handler http.Handler, address string) *Server {
	httpServer := &http.Server{
		Handler:           handler,
		Addr:              address,
		ReadTimeout:       _defaultReadTimeout,
		WriteTimeout:      _defaultWriteTimeout,
		IdleTimeout:       _defaultIdleTimeout,
		ReadHeaderTimeout: _defaultReadHeaderTimeout,
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
