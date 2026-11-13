package nullable

import (
	"bytes"
	"encoding/json"
	"time"
)

type Nullable[T any] struct {
	Value *T
	IsSet bool
}

func (o *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		o.IsSet = true
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.Value = &value
	o.IsSet = true
	return nil
}

type String = Nullable[string]
type Bool = Nullable[bool]
type Int = Nullable[int]
type Int8 = Nullable[int8]
type Int16 = Nullable[int16]
type Int32 = Nullable[int32]
type Int64 = Nullable[int64]
type Uint = Nullable[uint]
type Uint8 = Nullable[uint8]
type Uint16 = Nullable[uint16]
type Uint32 = Nullable[uint32]
type Uint64 = Nullable[uint64]
type Float32 = Nullable[float32]
type Float64 = Nullable[float64]
type Time = Nullable[time.Time]
