package domain

import (
	"bytes"
	"errors"

	"database/sql/driver"
)

// JSONB represents slice of bytes
type JSONB []byte

// Value impl.
func (j JSONB) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

// Scan impl.
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not string")
	}

	*j = append((*j)[0:0], s...)
	return nil
}

// MarshalJSON impl.
func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON impl.
func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// IsNull impl.
func (j JSONB) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// Equals impl.
func (j JSONB) Equals(j1 JSONB) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}
