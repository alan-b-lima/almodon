package middleware

import (
	"io"
	"log"

	"github.com/alan-b-lima/ansi-escape-sequences"
)

type Logger struct {
	il   log.Logger
	ansi bool
}

func NewLogger(w io.Writer, name string) *Logger {
	l := new(Logger)

	l.ansi = enableAnsi(w)

	l.il.SetOutput(w)
	l.il.SetFlags(log.Ldate | log.Ltime)
	l.il.SetPrefix(name + "> ")

	return l
}

func (l *Logger) Print(v ...any) {
	l.il.Print(v...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.il.Printf(format, v...)
}

func (l *Logger) Println(v ...any) {
	l.il.Println(v...)
}

func enableAnsi(w io.Writer) bool {
	f, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return false
	}

	err := ansi.EnableVirtualTerminal(f.Fd())
	return err == nil
}
