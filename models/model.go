package models

import (
	"fmt"

	countryCodes "github.com/launchdarkly/go-country-codes"
)

// Customer represents a customer and provides validation functionality
type Customer struct {
	ID             string    `json:"id"`
	Name           string    `json:"name" api:"required"`
	Contact        string    `json:"contact"` // optional field
	ContactNumber  string    `json:"contact_number" api:"required"`
	Address        Address   `json:"address" api:"required"`
	WeatherDetails []Weather `json:"weather"`
}

// Validate verifies data about the Customer. It does not duplicate verification of properties annotated with `api:"required"` tags
func (customer Customer) Validate() error {
	contactNumberMinLength := 7
	if len(customer.ContactNumber) < contactNumberMinLength {
		return fmt.Errorf("The contact number must be a minimum of %d digits", contactNumberMinLength)
	}
	return customer.Address.Validate()
}

// Customers is a list of Customer objects
type Customers []Customer

// Address represents a physical address
type Address struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"-"`
	// State           string          `json:"state"`
	// Address1        string          `json:"address_1"`
	// Address2        string          `json:"address_2"`
	// ZipCode         string          `json:"zip_code"`

	// For future we should support translating address fields to geo coordinates leveraging a geo API to be more flexible when obtaining address info from
	// the client
	// GeoCooordinates geo.Coordinates `json:"geo_coordinates"`
}

var errAddressValidationFailure = fmt.Errorf(`Please ensure city and country are provided`)

// Validate verifies that the city and country are specified as we'll need them to query the weather forecast API
func (address Address) Validate() error {
	if len(address.City) == 0 || len(address.Country) == 0 {
		return errAddressValidationFailure
	}

	_, err := address.SetCountryCode()
	return err
}

// SetCountryCode translates the Country specified to the code provided by Alpha2 codes standardized by ISO-3116
// https://www.iso.org/iso-3166-country-codes.html
func (address Address) SetCountryCode() (Address, error) {
	var matchedCode string
	matches := countryCodes.FindByName(address.Country)
	if len(matches) != 1 {
		countryCode, ok := countryCodes.GetByAlpha2(address.Country)
		if !ok {
			countryCode, ok = countryCodes.GetByAlpha3(address.Country)
			if !ok {
				return address, fmt.Errorf("Associated country could not be resolved")
			}
		}
		matchedCode = countryCode.Alpha2

	} else {
		matchedCode = matches[0].Alpha2
	}

	address.CountryCode = matchedCode

	return address, nil
}
