package router

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAPIAnnotation(t *testing.T) {
	type TestRequest struct {
		StringKey1  string   `json:"string_key_1" api:"required"`
		StringKey2  string   `json:"string_key_2"`
		FloatPtrKey *float64 `json:"float_ptr_key"`
	}
	var floatVar float64
	tests := []struct {
		name     string
		input    map[string]interface{}
		output   interface{}
		expError error
	}{
		{
			name:     "Output is non struct type",
			input:    map[string]interface{}{},
			output:   floatVar,
			expError: fmt.Errorf("Output object should be as struct or a pointer to one"),
		},
		{
			name:     "input is nil",
			input:    nil,
			output:   TestRequest{},
			expError: fmt.Errorf("string_key_1 required"),
		},
		{
			name:     "input is empty map",
			input:    map[string]interface{}{},
			output:   TestRequest{},
			expError: fmt.Errorf("string_key_1 required"),
		},
		{
			name:     "input has no required keys",
			input:    map[string]interface{}{"string_key_2": "hello"},
			output:   TestRequest{},
			expError: fmt.Errorf("string_key_1 required"),
		},
		{
			name:     "input has required keys all present",
			input:    map[string]interface{}{"string_key_1": "hello", "string_key_2": "world"},
			output:   TestRequest{},
			expError: nil,
		},
		{
			name:     "input has string required property empty",
			input:    map[string]interface{}{"string_key_1": "", "string_key_2": "world"},
			output:   TestRequest{},
			expError: fmt.Errorf("string_key_1 required"),
		},
		{
			name:     "output is ptr to struct",
			input:    map[string]interface{}{"string_key_1": "hello", "string_key_2": "world"},
			output:   &TestRequest{},
			expError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recErr := validateAPIAnnotation(test.input, test.output)
			assert.Equal(t, test.expError, recErr)
		})
	}
}
