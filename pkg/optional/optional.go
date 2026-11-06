package optional

import (
	"bytes"
	"encoding/json"
)

type String struct {
	Value *string
	IsSet bool
}

func (o *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		o.IsSet = true
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.Value = &s
	o.IsSet = true
	return nil
}
