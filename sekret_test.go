package sekret_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lex-unix/sekret"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSekret(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected string
	}

	tests := []test{
		{
			name:     "masks secret from string",
			input:    "openai api token",
			expected: "******",
		},
		{
			name:     "masks empty string",
			input:    "",
			expected: "******",
		},
		{
			name:     "masks single character",
			input:    "x",
			expected: "******",
		},
		{
			name:     "masks very long string",
			input:    "this-is-a-very-long-secret-password-with-many-characters-1234567890",
			expected: "******",
		},
		{
			name:     "masks string with special characters",
			input:    "p@ssw0rd!#$%^&*()",
			expected: "******",
		},
		{
			name:     "masks string with unicode",
			input:    "пароль123",
			expected: "******",
		},
		{
			name:     "masks string with newlines",
			input:    "line1\nline2\nline3",
			expected: "******",
		},
		{
			name:     "masks string with tabs and spaces",
			input:    "secret\t\twith\ttabs and spaces",
			expected: "******",
		},
		{
			name:     "masks numeric string",
			input:    "1234567890",
			expected: "******",
		},
		{
			name:     "masks json-like string",
			input:    `{"key": "secret-value"}`,
			expected: "******",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := sekret.New(tt.input)
			assert.Equal(t, tt.expected, secret.String())
			assert.Equal(t, tt.input, secret.ExposeSecret())
		})
	}
}

var jsonString = []byte(`{
	"username": "lex",
	"password": "secret-password",
	"db": {
		"host": "test.com",
		"password": "db-password"
	}
}`)

var jsonMarshalledString = []byte(`{
	"username": "lex",
	"password": "******",
	"db": {
		"host": "test.com",
		"password": "******"
	}
}`)

var yamlString = []byte(`
username: lex
password: secret-password
db:
  host: test.com
  password: db-password
`)

var yamlMarshalledString = []byte(`
username: lex
password: "******"
db:
  host: test.com
  password: "******"
`)

type dummy struct {
	Username string                `json:"username" yaml:"username"`
	Password sekret.Sekret[string] `json:"password" yaml:"password"`
	DB       struct {
		Host     string                `json:"host" yaml:"host"`
		Passwrod sekret.Sekret[string] `json:"password" yaml:"password"`
	} `json:"db" yaml:"db"`
}

func TestUmarshallJSON(t *testing.T) {
	var s dummy
	err := json.Unmarshal(jsonString, &s)
	assert.NoError(t, err)
	assert.Equal(t, "lex", s.Username)
	assert.Equal(t, "******", s.Password.String())
	assert.Equal(t, "secret-password", s.Password.ExposeSecret())
}

func TestMarsalJSON(t *testing.T) {
	s := dummy{
		Username: "lex",
		Password: sekret.New("secret-password"),
	}
	s.DB.Host = "test.com"
	s.DB.Passwrod = sekret.New("db-password")
	raw, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.JSONEq(t, string(jsonMarshalledString), string(raw))
}

func TestUnmarshallYAML(t *testing.T) {
	var s dummy
	err := yaml.Unmarshal(yamlString, &s)
	assert.NoError(t, err)
	assert.Equal(t, "lex", s.Username)
	assert.Equal(t, "******", s.Password.String())
	assert.Equal(t, "secret-password", s.Password.ExposeSecret())
}

func TestMarshallYAML(t *testing.T) {
	s := dummy{
		Username: "lex",
		Password: sekret.New("secret-password"),
	}
	s.DB.Host = "test.com"
	s.DB.Passwrod = sekret.New("db-password")
	raw, err := yaml.Marshal(s)
	assert.NoError(t, err)
	assert.YAMLEq(t, string(yamlMarshalledString), string(raw))
}

func TestSekretFormatting(t *testing.T) {
	secret := sekret.New("my-secret")

	t.Run("fmt.Sprintf with %s", func(t *testing.T) {
		result := fmt.Sprintf("Password: %s", secret)
		assert.Equal(t, "Password: ******", result)
	})

	t.Run("fmt.Sprintf with %v", func(t *testing.T) {
		result := fmt.Sprintf("Password: %v", secret)
		assert.Equal(t, "Password: ******", result)
	})

	t.Run("fmt.Sprintf with %+v", func(t *testing.T) {
		result := fmt.Sprintf("Password: %+v", secret)
		assert.Equal(t, "Password: ******", result)
	})

	t.Run("fmt.Sprintf with %#v", func(t *testing.T) {
		result := fmt.Sprintf("Password: %#v", secret)
		assert.Contains(t, result, "sekret.Sekret{******}")
	})
}
