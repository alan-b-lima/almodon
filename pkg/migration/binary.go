package migration

import (
	"encoding/binary"
	"errors"
	"slices"
)

func (m Migration) LenBinary() int {
	// sizeof(m.hash) + sizeof(len(m.scripts)) + sizeof(m.hash) * len(m.scripts)
	return len(Hash{}) + 8 + len(Hash{})*len(m.scripts)
}

func (m Migration) AppendBinary(buf []byte) ([]byte, error) {
	buf = slices.Grow(buf, m.LenBinary())

	buf = append(buf, m.hash[:]...)
	buf = binary.LittleEndian.AppendUint64(buf, uint64(len(m.scripts)))
	for _, hash := range m.scripts {
		buf = append(buf, hash[:]...)
	}

	return buf, nil
}

func (m Migration) MarshalBinary() ([]byte, error) {
	return m.AppendBinary(nil)
}

var ErrInvalidBinary = errors.New("migration: invalid binary encoding")

func (m *Migration) UnmarshalBinary(buf []byte) error {
	if len(buf) < len(Hash{})+8 {
		return ErrInvalidBinary
	}

	var hash Hash
	n := copy(hash[:], buf)
	buf = buf[n:]

	size := binary.LittleEndian.Uint64(buf)
	buf = buf[8:]

	if len(buf) < len(Hash{})*int(size) {
		return ErrInvalidBinary
	}

	scripts := make([]Hash, 0, size)
	for range size {
		var hash Hash
		n := copy(hash[:], buf)
		scripts = append(scripts, hash)
		buf = buf[n:]
	}

	*m = Migration{
		hash:    hash,
		scripts: scripts,
	}
	return nil
}
