package router

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Request represents the data associated with a handler
type Request struct {
	// Info represents the json query or method parameters associated with a handler
	Info map[string]interface{} `json:"info"`
}

// Parse deserializes the request object into the output param. It provides validation of the request based on "api" annotated properties
func (req Request) Parse(out interface{}) error {
	info, err := json.Marshal(req.Info)
	if err != nil {
		return fmt.Errorf("Request parameters couldn't be converted to json string. Err: %s", err.Error())
	}

	if err = json.Unmarshal(info, out); err != nil {
		return fmt.Errorf("Failed to unmarshal request parameters. Err: %s", err.Error())
	}

	if err = validateAPIAnnotation(req.Info, out); err != nil {
		return fmt.Errorf("Request validation failed: %s", err.Error())
	}

	return nil
}

// validateAPIAnnotation validates the `api:"required"` annotated tags in the request struct
func validateAPIAnnotation(info map[string]interface{}, req interface{}) error {
	kind := reflect.TypeOf(req).Kind()
	if kind == reflect.Ptr {
		return validateAPIAnnotation(info, reflect.Indirect(reflect.ValueOf(req)).Interface())
	} else if kind != reflect.Struct {
		return fmt.Errorf("Output object should be as struct or a pointer to one")
	}

	val := reflect.ValueOf(req)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.Tag.Get("api") != "required" {
			continue
		}

		jsonName := field.Tag.Get("json")
		fieldRequiredErr := fmt.Errorf("%s required", jsonName)
		if info == nil {
			return fieldRequiredErr
		}

		if _, ok := info[jsonName]; !ok {
			return fieldRequiredErr
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			subInfo := info[jsonName].(map[string]interface{})
			return validateAPIAnnotation(subInfo, val.Field(i).Interface())
		case reflect.String:
			strValue := info[jsonName].(string)
			if strValue == "" {
				return fieldRequiredErr
			}
		}
	}
	return nil
}
