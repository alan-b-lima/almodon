package request

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/uuid"
)

type Service interface {
	
}

type (
	Create struct {
		Author uuid.UUID `json:"-"`
		Title  string    `json:"title"`
		Memo   string    `json:"memo"`
	}
)

type (
	Result struct {
		UUID     uuid.UUID `json:"uuid"`
		Number   int       `json:"number"`
		Author   uuid.UUID `json:"author"`
		Title    string    `json:"title"`
		Memo     string    `json:"memo"`
		Status   Status    `json:"status"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
	}
)
