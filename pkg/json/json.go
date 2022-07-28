package json

import (
	"encoding/json"
	"os"
)

func WriteJSONFile(filename string, v any) error {
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

func ReadJSONFile[T any](filename string) (T, error) {
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
