package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GenerateMessage(err error) map[string]string {
	var vErr validator.ValidationErrors
	if errors.As(err, &vErr) {
		messages := make(map[string]string, len(vErr))
		for _, v := range vErr {
			switch v.Tag() {
			case "email":
				messages[v.Field()] = fmt.Sprintf("%s is not valid email", v.Value())
			case "required":
				messages[v.Field()] = fmt.Sprintf("%s is required", v.Field())
			case "min":
				messages[v.Field()] = fmt.Sprintf("%s must be at least %s characters", v.Field(), v.Param())
			case "max":
				messages[v.Field()] = fmt.Sprintf("%s must be at most %s characters", v.Field(), v.Param())
			case "len":
				messages[v.Field()] = fmt.Sprintf("%s must be exactly %s characters", v.Field(), v.Param())
			case "uuid":
				messages[v.Field()] = fmt.Sprintf("%s is not a valid UUID", v.Field())
			case "alpha":
				messages[v.Field()] = fmt.Sprintf("%s must only contain alphabetic characters", v.Field())
			case "alpha_dash":
				messages[v.Field()] = fmt.Sprintf("%s must only contain alphabetic characters, dashes, and underscores", v.Field())
			case "alpha_num":
				messages[v.Field()] = fmt.Sprintf("%s must only contain alphabetic characters and numbers", v.Field())
			case "numeric":
				messages[v.Field()] = fmt.Sprintf("%s must be a numeric value", v.Field())
			case "gt":
				messages[v.Field()] = fmt.Sprintf("%s must be greater than %s", v.Field(), v.Param())
			case "gte":
				messages[v.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", v.Field(), v.Param())
			case "lt":
				messages[v.Field()] = fmt.Sprintf("%s must be less than %s", v.Field(), v.Param())
			case "lte":
				messages[v.Field()] = fmt.Sprintf("%s must be less than or equal to %s", v.Field(), v.Param())
			case "url":
				messages[v.Field()] = fmt.Sprintf("%s is not a valid URL", v.Field())
			case "hex":
				messages[v.Field()] = fmt.Sprintf("%s must be a valid hexadecimal string", v.Field())
			case "date":
				messages[v.Field()] = fmt.Sprintf("%s is not a valid date, e.g: 2006-01-02", v.Field())
			case "timezone":
				messages[v.Field()] = fmt.Sprintf("%s is not a valid timezone e.g: UTC,+08:00,Asia,Jakarta,America,New_York", v.Field())
			case "ip":
				messages[v.Field()] = fmt.Sprintf("%s is not a valid IP address", v.Field())
			}
		}
		return messages
	}
	return map[string]string{"error": err.Error()}
}
