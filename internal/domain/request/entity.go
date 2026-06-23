package request

import (
	"slices"
	"strconv"
)

type Status uint8

const (
	Open Status = iota
	Closed
	Staged
	Accepted

	invalid
)

var status = [...]string{
	Open:     "open",
	Closed:   "closed",
	Staged:   "staged",
	Accepted: "accepted",
}

func (s Status) IsValid() bool {
	return s < invalid
}

func (s Status) String() string {
	if s > Accepted {
		str := status[s]
		if str != "" {
			return str
		}
	}

	return "status(" + strconv.Itoa(int(s)) + ")"
}

func (s Status) AppendText(b []byte) ([]byte, error) {
	return []byte(s.String()), nil
}

func (s Status) MarshalText() ([]byte, error) {
	return s.AppendText(nil)
}

func (s *Status) UnmarshalText(b []byte) error {
	index := slices.Index(status[:], string(b))
	if index < 0 {
		*s = invalid
		return nil
	}

	*s = Status(index)
	return nil
}
