package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alan-b-lima/almodon/internal/api/v1"
	"github.com/alan-b-lima/almodon/internal/middleware"

	"github.com/alan-b-lima/ansi-escape-sequences"
)

var StdOut = os.Stdout

func main() {
	log := middleware.NewLogger(StdOut, "")
	style := Styles()

	ln, err := net.Listen("tcp", ":4545")
	if err != nil {
		log.Println(err)
		return
	}

	url := "http://" + strings.Replace(ln.Addr().String(), "[::]", "localhost", 1)
	log.Printf("Server listening at %s\n", style.HyperLink(url))

	srv := http.Server{Handler: TrafficLogMiddleware(log, style, api.New())}
	done := EnableGracefulShutdown(func() {
		log.Println("Closing...")
		srv.Shutdown(context.Background())
	})

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}

	<-done
}

func TrafficLogMiddleware(log *middleware.Logger, s Style, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.NewW(w)

		handler.ServeHTTP(rw, r)

		var pen ansi.Pen
		switch rw.StatusCode() / 100 {
		case 5: pen = s.ServerError
		case 4: pen = s.ClientError
		case 3: pen = s.Redirect
		case 2: pen = s.Success

		default:
			pen.SetStyle(false)
		}

		var b strings.Builder
		pen.Writer = &b

		io.WriteString(&pen, fmt.Sprintf(" %03d ", rw.StatusCode()))
		log.Printf("%s %s %s %s\n", b.String(), r.RemoteAddr, r.Method, r.URL)
	}
}

func EnableGracefulShutdown(fn func()) <-chan struct{} {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan struct{}, 1)

	go func() {
		<-signals
		fn()
		done <- struct{}{}
	}()

	return done
}

type Style struct {
	HyperLink func(string) string

	ServerError ansi.Pen
	ClientError ansi.Pen
	Redirect    ansi.Pen
	Success     ansi.Pen
	NoStyle     ansi.Pen

	Enabled bool
}

func Styles() (s Style) {
	s.Enabled = ansi.EnableVirtualTerminal(os.Stdout.Fd()) == nil

	if s.Enabled {
		s.Success.BGColor(ansi.RGBFromHex(0x0ed145))
		s.Success.FGColor(ansi.RGBFromHex(0xffffff))

		s.Redirect.BGColor(ansi.RGBFromHex(0x4b53cc))
		s.Redirect.FGColor(ansi.RGBFromHex(0xffffff))

		s.ClientError.BGColor(ansi.RGBFromHex(0xea1d1d))
		s.ClientError.FGColor(ansi.RGBFromHex(0xffffff))

		s.ServerError.BGColor(ansi.RGBFromHex(0x88001b))
		s.ServerError.FGColor(ansi.RGBFromHex(0xffffff))

		s.HyperLink = hyperlink

		return s
	}

	s.Success.SetStyle(false)
	s.Redirect.SetStyle(false)
	s.ClientError.SetStyle(false)
	s.ServerError.SetStyle(false)
	s.NoStyle.SetStyle(false)
	s.HyperLink = func(s string) string { return s }

	return s
}

func hyperlink(link string) string {
	var pen ansi.Pen
	pen.FGColor(ansi.RGBFromHex(0x4e8597))

	return pen.Sprint(ansi.HyperLinkP(link))
}
