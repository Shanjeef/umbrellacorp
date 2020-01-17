package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCustomer(t *testing.T) {
	tests := []struct {
		name   string
		input  Customer
		expErr error
	}{
		{
			name: "phone number empty",
			input: Customer{
				ContactNumber: "",
			},
			expErr: fmt.Errorf("The contact number must be a minimum of 7 digits"),
		},
		{
			name: "phone number less than minimum",
			input: Customer{
				ContactNumber: "3848",
			},
			expErr: fmt.Errorf("The contact number must be a minimum of 7 digits"),
		},
		{
			name: "City only provided",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					City: "Toronto",
				},
			},
			expErr: errAddressValidationFailure,
		},
		{
			name: "Country only provided",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					Country: "CA",
				},
			},
			expErr: errAddressValidationFailure,
		},
		{
			name: "City and Country Provided",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					City:    "Toronto",
					Country: "Canada",
				},
			},
			expErr: nil,
		},
		{
			name: "Country Unknown",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					City:    "Toronto",
					Country: "Fake Country",
				},
			},
			expErr: fmt.Errorf("Associated country could not be resolved"),
		},
		{
			name: "CA Location",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					City:    "Toronto",
					Country: "CA",
				},
			},
			expErr: nil,
		},
		{
			name: "US Location",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					City:    "Chicago",
					Country: "US",
				},
			},
			expErr: nil,
		},
		{
			name: "USA Location",
			input: Customer{
				ContactNumber: "4165555555",
				Address: Address{
					City:    "Chicago",
					Country: "USA",
				},
			},
			expErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recErr := test.input.Validate()
			assert.Equal(t, test.expErr, recErr)
		})
	}
}
