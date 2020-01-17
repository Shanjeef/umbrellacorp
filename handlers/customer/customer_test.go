package customer

import (
	"fmt"
	"os"
	"testing"
	"time"
	"umbrellacorp/components/weatherforecaster"
	"umbrellacorp/models"
	"umbrellacorp/router"
	"umbrellacorp/util"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.M) {
	weatherforecaster.Configure(true)
	os.Exit(t.Run())
}

func TestValidateUniqueCustomer(t *testing.T) {
	tests := []struct {
		Name              string
		existingCustomers models.Customers
		newCustomer       models.Customer
		expError          error
	}{
		{
			Name:              "no existing customers",
			existingCustomers: nil,
			newCustomer: models.Customer{
				Name:          "Awesome Company",
				ContactNumber: "5165555555",
			},
			expError: nil,
		},
		{
			Name: "customer with existing name",
			existingCustomers: models.Customers{
				{
					Name:          "Awesome Company",
					ContactNumber: "5165555555",
				},
			},
			newCustomer: models.Customer{
				Name:          "Awesome Company",
				ContactNumber: "5165555558",
			},
			expError: fmt.Errorf("An existing customer with the same name exists"),
		},
		{
			Name: "customer with existing contact number exists",
			existingCustomers: models.Customers{
				{
					Name:          "Fortune 500 Company",
					ContactNumber: "5165555555",
				},
			},
			newCustomer: models.Customer{
				Name:          "Awesome Company",
				ContactNumber: "5165555555",
			},
			expError: fmt.Errorf("An existing customer with the same contact number exists"),
		},
		{
			Name: "net new customer",
			existingCustomers: models.Customers{
				{
					Name:          "Fortune 500 Company",
					ContactNumber: "5165555555",
				},
			},
			newCustomer: models.Customer{
				Name:          "Awesome Company",
				ContactNumber: "5165555557",
			},
			expError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			recErr := validateUniqueCustomer(test.existingCustomers, test.newCustomer)
			assert.Equal(t, test.expError, recErr)
		})
	}
}

func TestSetCustomer(t *testing.T) {
	// Fetch exp weather details
	now := time.Date(2017, 02, 16, 0, 0, 0, 0, time.UTC)
	dateRange := util.DateRange{Start: now, End: now.AddDate(0, 0, 5)}
	weatherDetails, err := weatherforecaster.NewForecaster().UpcomingWeather("Toronto", "CA", dateRange, models.WeatherTypeRain)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tests := []struct {
		Name              string
		input             map[string]interface{}
		existingCustomers models.Customers
		expCustomers      models.Customers
		expError          error
	}{
		{
			Name: "create new customer",
			input: map[string]interface{}{
				"name":           "Awesome Company",
				"contact_number": "4165555555",
				"address": map[string]interface{}{
					"city":    "Toronto",
					"country": "CA",
				},
			},
			expCustomers: models.Customers{
				models.Customer{
					Name:          "Awesome Company",
					ContactNumber: "4165555555",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
					WeatherDetails: weatherDetails,
				},
			},
			expError: nil,
		},
		{
			Name: "Add a new customer that is already present",
			input: map[string]interface{}{
				"name":           "New Customer",
				"contact_number": "4165555555",
				"address": map[string]interface{}{
					"city":    "Toronto",
					"country": "CA",
				},
			},
			existingCustomers: models.Customers{
				models.Customer{
					ID:            "1",
					Name:          "Awesome Company",
					ContactNumber: "4165555555",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
				},
			},
			expCustomers: models.Customers{
				models.Customer{
					ID:            "1",
					Name:          "Awesome Company",
					ContactNumber: "4165555555",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
				},
			},
			expError: fmt.Errorf("An existing customer with the same contact number exists"),
		},
		{
			Name: "Update a customer that isn't present",
			input: map[string]interface{}{
				"ID":             "2", // Unknown customer ID
				"name":           "Awesome Company",
				"contact_number": "4165555555",
				"address": map[string]interface{}{
					"city":    "Toronto",
					"country": "CA",
				},
			},
			existingCustomers: models.Customers{
				models.Customer{
					ID:            "1",
					Name:          "Awesome Company",
					ContactNumber: "4165555555",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
				},
			},
			expCustomers: models.Customers{
				models.Customer{
					ID:            "1",
					Name:          "Awesome Company",
					ContactNumber: "4165555555",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
				},
			},
			expError: fmt.Errorf("Failed to locate existing customer with id: 2"),
		},
		{
			Name: "Update a customer record successfully",
			input: map[string]interface{}{
				"ID":             "1",
				"name":           "Awesome Company",
				"contact_number": "4169999999", // Change of contact number
				"address": map[string]interface{}{
					"city":    "Toronto",
					"country": "CA",
				},
			},
			existingCustomers: models.Customers{
				models.Customer{
					ID:            "1",
					Name:          "Awesome Company",
					ContactNumber: "4165555555",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
				},
			},
			expCustomers: models.Customers{
				models.Customer{
					ID:            "1",
					Name:          "Awesome Company",
					ContactNumber: "4169999999",
					Address: models.Address{
						City:        "Toronto",
						Country:     "CA",
						CountryCode: "CA",
					},
					WeatherDetails: weatherDetails,
				},
			},
			expError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			customers = test.existingCustomers

			req := router.Request{Info: test.input}
			_, recErr := setCustomer(req)
			assert.Equal(t, test.expError, recErr)

			if len(test.expCustomers) != len(customers) {
				t.Fatalf("Exp customer size: %d, actual customers size: %d", len(test.expCustomers), len(customers))
			}

			// Copy the ID so that its effectively not compared
			for i := range test.expCustomers {
				assert.NotEmpty(t, customers[i].ID)
				test.expCustomers[i].ID = customers[i].ID
			}
			assert.Equal(t, test.expCustomers, customers)
		})
	}
}
