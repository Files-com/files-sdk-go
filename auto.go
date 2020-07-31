package files_sdk

import (
	"encoding/json"
)

type Auto struct {
	Dynamic json.RawMessage `json:"dynamic,omitempty"`
}

type AutoCollection []Auto

func (a *Auto) UnmarshalJSON(data []byte) error {
	type auto Auto
	var v auto
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Auto(v)
	return nil
}

func (a *AutoCollection) UnmarshalJSON(data []byte) error {
	type autos []Auto
	var v autos
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = AutoCollection(v)
	return nil
}
