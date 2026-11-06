package optional

import "encoding/json"

type String struct {
	Value *string
	IsSet bool
}

func (o *String) UnmarshalJSON(data []byte) error {
	o.IsSet = true

	if string(data) == "null" {
		o.Value = nil
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.Value = &s
	return nil
}
