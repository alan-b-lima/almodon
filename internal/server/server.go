package server

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alan-b-lima/almodon/internal/client"
	sessions "github.com/alan-b-lima/almodon/internal/domain/session/resource"
	"github.com/alan-b-lima/almodon/internal/support/resource"
)

func NewAlmodon() (http.Handler, *API, error) {
	client, err := client.New()
	if err != nil {
		return nil, nil, err
	}

	// Mount the API last not to handle conditinal
	// closing, see [domain.New] for an example.
	api, err := NewAPI()
	if err != nil {
		return nil, nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/", client)

	return sessions.Wrap(mux), api, nil
}

type Server struct {
	Listener net.Listener
	Handler  http.Handler
	CertFile string
	KeyFile  string
	Logger   *slog.Logger

	server http.Server
	base   context.Context
	close  context.CancelFunc
}

func (s *Server) Serve() error {
	if s.Handler == nil {
		panic("nil handler")
	}

	var olog *log.Logger
	if s.Logger != nil {
		olog = slog.NewLogLogger(s.Logger.Handler(), slog.LevelError)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	s.base, s.close = context.WithCancel(context.Background())

	s.server = http.Server{
		Handler:  logger(s.Logger, s.Handler),
		ErrorLog: olog,
		BaseContext: func(net.Listener) context.Context {
			return s.base
		},
	}

	go s.parseSignals(signals)

	authority := strings.Replace(s.Addr().String(), "[::]", "localhost", 1)

	if s.CertFile == "" {
		s.Logger.Info("server listening at http://" + authority)
		return s.server.Serve(s.Listener)
	} else {
		s.Logger.Info("server listening at https://" + authority)
		return s.server.ServeTLS(s.Listener, s.CertFile, s.KeyFile)
	}
}

func (s *Server) Addr() net.Addr {
	return s.Listener.Addr()
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer s.close()
	return s.server.Shutdown(ctx)
}

func (s *Server) ForceShutdown() error {
	defer s.close()
	return s.server.Close()
}

func (s *Server) parseSignals(signals chan os.Signal) {
	defer signal.Stop(signals)

	select {
	case <-s.base.Done():
		return

	case <-signals:
		go func() {
			if err := s.Shutdown(context.Background()); err != nil {
				s.Logger.Error("shutdown", slog.String("error", err.Error()))
			}
		}()
	}

	select {
	case <-s.base.Done():
		return

	case <-signals:
		s.ForceShutdown()
	}
}

func logger(log *slog.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := resource.NewResponseWriter(w)
		start := time.Now()

		h.ServeHTTP(rw, r)

		elapsed := time.Since(start)

		level := slog.LevelInfo
		if rw.StatusCode()/100 == 5 {
			level = slog.LevelError
		}

		log.LogAttrs(
			r.Context(), level, "completed request",
			slog.Int("status", rw.StatusCode()),
			slog.String("remote", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.String()),
			slog.Duration("elapsed", elapsed),
			slog.Any("error", rw.Error()),
		)
	}
}
