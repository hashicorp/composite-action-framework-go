package json

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/hashicorp/composite-action-framework-go/pkg/fs"
)

func Read[T any](r io.Reader) (T, error) {
	v := new(T)
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	err := d.Decode(v)
	return *v, err
}

func String(v any) (string, error) {
	buf := &bytes.Buffer{}
	err := Write(buf, v)
	return buf.String(), err
}

func Write(w io.Writer, v any) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(v)
}

func WriteFile(filename string, v any, fsOpts ...fs.Option) error {
	f, err := fs.Create(filename, fsOpts...)
	if err != nil {
		return err
	}
	var closeErr error
	defer func() { closeErr = f.Close() }()
	if err := Write(f, v); err != nil {
		return err
	}
	return closeErr
}

func ReadFile[T any](filename string) (T, error) {
	f, err := os.Open(filename)
	if err != nil {
		return *(new(T)), err
	}
	var closeErr error
	defer func() { closeErr = f.Close() }()
	v, err := Read[T](f)
	if err != nil {
		return v, err
	}
	return v, closeErr
}

func ReadBytes[T any](jsonBytes []byte) (T, error) {
	return Read[T](bytes.NewBuffer(jsonBytes))
}

func ReadString[T any](jsonString string) (T, error) {
	return ReadBytes[T]([]byte(jsonString))
}
