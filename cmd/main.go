package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"github.com/alan-b-lima/almodon/internal/server"
)

func main() {
	// os.Chdir("..")

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

	var httpHandler, httpsHandler http.ServeMux

	var (
		httpServer  = server.Server{Addr: ":4545", Handler: &httpHandler, Logger: logger}
		httpsServer = server.Server{Addr: ":4343", Handler: &httpsHandler, Logger: logger, CertFile: "server.rsa.crt", KeyFile: "server.rsa.key"}
	)

	httpHandler.Handle("/api/", http.RedirectHandler("https://"+httpServer.Addr, http.StatusMovedPermanently))
	httpHandler.Handle("/", handler)

	httpsHandler.Handle("/api/", http.NotFoundHandler())
	httpsHandler.Handle("/", handler)

	done := make(chan struct{}, 2)

	go func() {
		if err := httpServer.Serve(); err != nil {
			fmt.Println("http", err)
		}
		done <- struct{}{}
	}()

	go func() {
		if err := httpsServer.Serve(); err != nil {
			fmt.Println("https", err)
		}
		done <- struct{}{}
	}()

	<-done
	<-done
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
