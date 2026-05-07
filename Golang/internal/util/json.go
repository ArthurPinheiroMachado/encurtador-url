package util

import (
	"encoding/json"
	"io"
)

func JsonDecodeFromBytes[T any](bt []byte) (*T, error) {
	trace := CreateErrorContext("util.JSONDecodeFromBytes")

	var value T

	_err := trace.Apply(json.Unmarshal(bt, &value))

	if _err != nil {
		return nil, _err
	}

	return &value, nil
}

func JsonDecodeFromReader[T any](reader io.Reader) (*T, error) {
	trace := CreateErrorContext("util.JSONDecodeFromReader")

	var value T

	if err := json.NewDecoder(reader).Decode(&value); err != nil {
		return nil, trace.Apply(err)
	}

	return &value, nil
}

func JsonEncodeToBytes[T any](value T) ([]byte, error) {
	trace := CreateErrorContext("util.JSONEncodeToBytes")

	data, err := json.Marshal(value)
	if err != nil {
		return nil, trace.Apply(err)
	}

	return data, nil
}

func JsonEncodeToWriter[T any](writer io.Writer, value T) error {
	trace := CreateErrorContext("util.JSONEncodeToWriter")

	if err := json.NewEncoder(writer).Encode(value); err != nil {
		return trace.Apply(err)
	}

	return nil
}
