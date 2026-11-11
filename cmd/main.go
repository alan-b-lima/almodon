package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/alan-b-lima/almodon/internal/server"
)

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "version":
			fmt.Println(version)
			return

		case "copyright", "legal":
			fmt.Println(legal)
			return

		case "help":
			fmt.Println(help)
			return
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	handler, api, err := server.NewAlmodon()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer api.Close()

	addr := ":4545"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	var (
		lnPlain, chPlain = newSink(ln.Addr())
		lnTLS, chTLS     = newSink(ln.Addr())
	)

	go multiplex(context.Background(), ln, chPlain, chTLS)

	var httpHandler, httpsHandler http.ServeMux

	httpHandler.HandleFunc("/api/", redirectToHTTPS)
	httpHandler.Handle("/", handler)

	httpsHandler.Handle("/api/", handler)
	httpsHandler.Handle("/", http.NotFoundHandler())

	var wg sync.WaitGroup

	wg.Go(func() {
		server := server.Server{
			Listener: lnPlain,
			Handler:  &httpHandler,
			Logger:   logger.With(slog.String("scheme", "http")),
		}

		if err := server.Serve(); err != nil && err != http.ErrServerClosed {
			fmt.Println("http", err)
		}
	})

	wg.Go(func() {
		server := server.Server{
			Listener: lnTLS,
			Handler:  &httpsHandler,
			CertFile: "server.crt.pem",
			KeyFile:  "server.key.pem",
			Logger:   logger.With(slog.String("scheme", "https")),
		}

		if err := server.Serve(); err != nil && err != http.ErrServerClosed {
			fmt.Println("https", err)
		}
	})

	wg.Wait()
}

func multiplex(ctx context.Context, ln net.Listener, plain, tls chan<- net.Conn) error {
	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}

			continue
		}

		go func() {
			var b [1]byte
			n, err := conn.Read(b[:])
			if err != nil || n < len(b) {
				conn.Close()
				return
			}

			if b[0] == 0x16 {
				tls <- newConn(b[:], conn)
			} else {
				plain <- newConn(b[:], conn)
			}
		}()
	}
}

func newSink(addr net.Addr) (net.Listener, chan<- net.Conn) {
	conns := make(chan net.Conn)
	return listener{source: conns, laddr: addr}, conns
}

type listener struct {
	source chan net.Conn
	laddr  net.Addr
}

func (l listener) Accept() (net.Conn, error) {
	conn, ok := <-l.source
	if !ok {
		return nil, net.ErrClosed
	}

	return conn, nil
}

func (l listener) Addr() net.Addr {
	return l.laddr
}

func (l listener) Close() error {
	close(l.source)
	return nil
}

func newConn(b []byte, c net.Conn) net.Conn {
	return conn{
		Conn: c,
		mr:   io.MultiReader(bytes.NewReader(b), c),
	}
}

type conn struct {
	net.Conn
	mr io.Reader
}

func (c conn) Read(p []byte) (int, error) {
	return c.mr.Read(p)
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	url := "https://" + r.Host + r.RequestURI
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

const (
	version = `Almodon ` + tag + ` ` + runtime.GOOS + `/` + runtime.GOARCH
	legal   = version + copyright
	help    = version + "\n\n" + `Work in Progress`
)

const tag = `v0.0.1`

const copyright = `
Copyright (C) 2026 Alan Lima

Esse programa está licenciado sob a Licença GPLv3 (GNU General Public License
versão 3). Este programa é distribuído na esperança de ser útil, mas SEM
NENHUMA GARANTIA; sem mesmo a garantia implícita de COMERCIALIZAÇÃO ou
ADEQUAÇÃO A UM PROPÓSITO ESPECÍFICO. Consulte a Licença GPLv3 para mais
detalhes.

Você deve ter recebido uma cópia da Licença GPLv3 junto com este programa. Se
não, veja <https://www.gnu.org/licenses/>.`
