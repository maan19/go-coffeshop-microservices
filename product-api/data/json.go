package data

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(i)
}

// FromJSON deserializes a string based JSON into an object of type interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
