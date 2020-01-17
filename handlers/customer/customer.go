package customer

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"umbrellacorp/components/weatherforecaster"
	"umbrellacorp/models"
	"umbrellacorp/router"
	"umbrellacorp/util"
)

// Init registers handlers with the router
func Init() {
	routes := router.Routes{
		{
			Name:        "Get Customers",
			Methods:     []string{http.MethodGet},
			Path:        "/customers",
			HandlerFunc: getCustomers,
		},
		{
			Name:        "Set Customer",
			Methods:     []string{http.MethodPost, http.MethodPut},
			Path:        "/customers",
			HandlerFunc: setCustomer,
		},
	}
	router.RegisterRoutes("customer", routes)
}

var customers models.Customers

func getCustomers(req router.Request) (router.Response, error) {
	resp := router.Response{Info: map[string]interface{}{}}

	// TODO: Customers' weather forecast should be accurate. One option would be to fetch weather details here but we shouldn't
	// couple the client's request with 3rd party here. A better option would be to have an async task on our server that updates customers' weather details
	resp.Info["customers"] = customers
	return resp, nil
}

// setCustomer upserts a customer entry
func setCustomer(req router.Request) (router.Response, error) {
	resp := router.Response{Info: map[string]interface{}{}}
	var customer models.Customer
	err := req.Parse(&customer)
	if err != nil {
		return resp, err
	}

	if err := customer.Validate(); err != nil {
		return resp, err
	}

	addressModified := false
	prevAddress := customer.Address
	customer.Address, err = customer.Address.SetCountryCode()
	if err != nil {
		return resp, err
	}

	if !reflect.DeepEqual(prevAddress, customer.Address) {
		addressModified = true
	}

	updateWeatherDetailsFn := func(cus models.Customer) (models.Customer, error) {
		if !addressModified {
			return cus, nil
		}

		weatherDetails, err := fetchForecast(cus.Address)
		if err != nil {
			return cus, fmt.Errorf("Failed to obtain upcoming weather: %s", err.Error())
		}
		cus.WeatherDetails = weatherDetails
		return cus, nil
	}

	if customer.ID != "" {
		customer, err = updateWeatherDetailsFn(customer)
		if err != nil {
			return resp, err
		}

		if err := updateCustomer(customers, customer); err != nil {
			return resp, err
		}

	} else {
		if err := validateUniqueCustomer(customers, customer); err != nil {
			return resp, err
		}
		customer.ID = util.NewID()

		customer, err = updateWeatherDetailsFn(customer)
		if err != nil {
			return resp, err
		}

		customers = append(customers, customer)
	}

	resp.Info["customer"] = customer
	return resp, err
}

// validateUniqueCustomer verifies that there isn't an existing customer with the same name or contact number
func validateUniqueCustomer(existingCustomers models.Customers, customer models.Customer) error {
	for _, existingCustomer := range existingCustomers {
		if strings.ToLower(existingCustomer.Name) == strings.ToLower(customer.Name) {
			return fmt.Errorf("An existing customer with the same name exists")
		}

		if existingCustomer.ContactNumber == customer.ContactNumber {
			return fmt.Errorf("An existing customer with the same contact number exists")
		}
	}
	return nil
}

func updateCustomer(existingCustomers models.Customers, customer models.Customer) error {
	located := false
	for i, existingCustomer := range existingCustomers {
		if existingCustomer.ID == customer.ID {
			existingCustomers[i] = customer
			located = true
			break
		}
	}
	if !located {
		return fmt.Errorf("Failed to locate existing customer with id: %s", customer.ID)
	}
	return nil
}

func fetchForecast(address models.Address) ([]models.Weather, error) {
	// Forecaster API seems to only give data from Feb 2017
	now := time.Date(2017, 02, 16, 0, 0, 0, 0, time.UTC)
	dateRange := util.DateRange{Start: now, End: now.AddDate(0, 0, 5)}
	return weatherforecaster.NewForecaster().UpcomingWeather(address.City, address.CountryCode, dateRange, models.WeatherTypeRain)
}
