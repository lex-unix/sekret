// Package sekret provides types for handling sensitive data safely.
// Secrets are automatically masked when printed with fmt functions, JSON, and YAML serialization.
package sekret

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

const mask = "******"

// Sekret wraps sensitive values and prevents accidental exposure in fmt output.
type Sekret[T any] struct {
	val T
}

// New creates a new Sekret wrapping the given value.
func New[T any](val T) Sekret[T] {
	return Sekret[T]{val: val}
}

// ExposeSecret returns the actual secret value. Use deliberately when the raw value is needed.
func (s Sekret[T]) ExposeSecret() T {
	return s.val
}

// String returns a masked string for fmt.Printf, fmt.Sprintf, and similar operations.
func (s Sekret[T]) String() string {
	return mask
}

// GoString returns a masked string for fmt %#v formatting.
func (s Sekret[T]) GoString() string {
	return "sekret.Sekret{******}"
}

// MarshalJSON masks the secret when encoding to JSON.
func (s Sekret[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(mask)
}

// UnmarshalJSON loads a secret from JSON data.
func (s *Sekret[T]) UnmarshalJSON(data []byte) error {
	var val T
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	s.val = val
	return nil
}

// MarshalYAML masks the secret when encoding to YAML.
func (s Sekret[T]) MarshalYAML() (any, error) {
	return mask, nil
}

// UnmarshalYAML loads a secret from YAML data.
func (s *Sekret[T]) UnmarshalYAML(value *yaml.Node) error {
	var temp T
	if err := value.Decode(&temp); err != nil {
		return err
	}
	s.val = temp
	return nil
}
