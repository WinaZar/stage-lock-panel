package utils

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
)

type customValidator struct{}

func (cv *customValidator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

// ParseValidatorErrors func parse govalidator multiple errors to slice
func ParseValidatorErrors(err error) map[string][]string {
	result := make(map[string][]string)
	errs := err.(govalidator.Errors).Errors()
	for _, e := range errs {
		split := strings.Split(e.Error(), ":")
		field, message := split[0], strings.Trim(split[1], " ")
		value, exists := result[field]
		if exists == true {
			result[field] = append(value, message)
		} else {
			value = make([]string, 0)
			result[field] = append(value, message)

		}
	}
	return result
}

// GetAndValidateLockData func bind lockData and validate it
func GetAndValidateLockData(ctx *CustomContext, lockData *LockData) map[string][]string {
	result := make(map[string][]string)

	if err := ctx.Bind(lockData); err != nil {
		result["bind"] = []string{err.Error()}
	} else if err := ctx.Validate(lockData); err != nil {
		fmt.Println(lockData.Code)
		result = ParseValidatorErrors(err)
	}
	return result
}

// Validator is a customValidator instance
var Validator = &customValidator{}
