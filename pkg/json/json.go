package json

import (
	"bytes"
	"encoding/json"
	"os"
)

func WriteFile(filename string, v any) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	var closeErr error
	defer func() { closeErr = f.Close() }()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		return err
	}
	return closeErr
}

func ReadFile[T any](filename string) (T, error) {
	v := new(T)
	f, err := os.Open(filename)
	if err != nil {
		return *v, err
	}
	var closeErr error
	defer func() { closeErr = f.Close() }()

	if err := json.NewDecoder(f).Decode(v); err != nil {
		return *v, err
	}

	return *v, closeErr
}

func ReadBytes[T any](jsonBytes []byte) (T, error) {
	v := new(T)
	buf := bytes.NewBuffer(jsonBytes)
	if err := json.NewDecoder(buf).Decode(v); err != nil {
		return *v, err
	}
	return *v, nil
}

func ReadString[T any](jsonString string) (T, error) {
	return ReadBytes[T]([]byte(jsonString))
}
