package kuanzhan

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
)

func jsonToMap(v any, m *map[string]interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	return dec.Decode(m)
}

func getError(resp any) error {
	s := structs.New(resp)
	if s.Field("Code").Value() != 200 {
		return errors.New(s.Field("Msg").Value().(string))
	}
	return nil
}
