package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type Version struct {
	Major   int
	Minor   int
	Patch   int
	MetaTag string
}

var errMalformedVersion = errors.New(`expected format <major>.<minor>.<patch>-<metatag>`)

func V(major, minor, patch int) Version {
	return Version{Major: major, Minor: minor, Patch: patch}
}

func Vm(major, minor, patch int, meta string) Version {
	return Version{Major: major, Minor: minor, Patch: patch, MetaTag: meta}
}

func (v Version) AppendText(b []byte) ([]byte, error) {
	if v.MetaTag == "" {
		b = fmt.Appendf(b, "%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else {
		b = fmt.Appendf(b, "%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.MetaTag)
	}

	return b, nil
}

func (v Version) MarshalText() ([]byte, error) {
	return v.AppendText(nil)
}

func (v *Version) UnmarshalText(b []byte) error {
	var ver Version
	r := bytes.NewReader(b)

	n, err := fmt.Fscanf(r, "%d.%d.%d", &ver.Major, &ver.Minor, &ver.Patch)
	if err != nil {
		return err
	}
	if n != 3 {
		return errMalformedVersion
	}

	byte, err := r.ReadByte()
	if err == io.EOF {
		*v = ver
		return nil
	}
	if byte != '-' {
		return errMalformedVersion
	}

	*v = ver
	v.MetaTag = string(b[len(b)-r.Len():])

	return nil
}
