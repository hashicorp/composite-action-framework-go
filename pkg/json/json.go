package json

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

func Read[T any](r io.Reader) (T, error) {
	v := new(T)
	err := json.NewDecoder(r).Decode(v)
	return *v, err
}

func Write(w io.Writer, v any) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(v)
}

func WriteFile(filename string, v any) error {
	f, err := os.Create(filename)
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
