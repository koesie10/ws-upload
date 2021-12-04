package wsupload

import (
	"encoding/json"
)

type Nullable interface {
	IsNull() bool
	Value() interface{}
}

var _ Nullable = NullFloat64{}
var _ Nullable = NullInt64{}

var _ json.Marshaler = NullFloat64{}
var _ json.Unmarshaler = &NullFloat64{}

type NullFloat64 struct {
	Valid   bool
	Float64 float64
}

func (v NullFloat64) IsNull() bool {
	return !v.Valid
}

func (v NullFloat64) Value() interface{} {
	return v.Float64
}

// MarshalJSON is the implementation of json.Marshaler
func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON is the implementation of json.Unmarshaler
func (v *NullFloat64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

var _ json.Marshaler = NullInt64{}
var _ json.Unmarshaler = &NullInt64{}

type NullInt64 struct {
	Valid bool
	Int64 int64
}

func (v NullInt64) IsNull() bool {
	return !v.Valid
}

func (v NullInt64) Value() interface{} {
	return v.Int64
}

// MarshalJSON is the implementation of json.Marshaler
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON is the implementation of json.Unmarshaler
func (v *NullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}
