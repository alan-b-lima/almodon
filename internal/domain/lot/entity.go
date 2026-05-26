package lot

import (
	"time"

	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/pkg/uuid"
)

const (
	SuplierMaxLen = 128
	NoteMaxLen    = 2048
)

func ProcessSupplier(supplier string) (string, error) {
	if !service.StringLenMax(supplier, SuplierMaxLen) {
		return "", ErrSupplierTooLong
	}

	return supplier, nil
}

func ProcessArrival(arrival time.Time) (time.Time, error) {
	return arrival, nil
}

func ProcessNote(note string) (string, error) {
	if !service.StringLenMax(note, SuplierMaxLen) {
		return "", ErrNoteTooLong
	}

	return note, nil
}

func StatusState(order uuid.UUID) State {
	if order.IsNil() {
		return Open
	}

	return Closed
}

type State bool

const (
	Open   State = true
	Closed State = false
)
