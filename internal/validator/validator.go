package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the FieldErrors is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map

func (v *Validator) AddFieldError(key, message string) {
	// Check if the FieldErrors has been initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckFied() adds an error only if the validator fails

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank()returns true if a field is not empty
func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if the value has less than a specific ammount of chars
func (v *Validator) MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermitedInt() returns true if the value is in a list of valid ints
func (v *Validator) PermitedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
