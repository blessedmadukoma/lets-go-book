package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

// EmailRX allows us perform sanity checks using regex, and stores the compiled *regexp.Regexp which is more performant that re-parsing the pattern each time we need it
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Valid checks if there are any errors (i.e. if FieldErrors and NonFieldErrors are empty)
func (v *Validator) Valid() bool {
	return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

// AddFieldError adds a form field error
func (v *Validator) AddFieldError(key, message string) {
	// if map wasn't initialized, initialize
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// AddNonFieldError adds a non form field error to the slice
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField checks if a field is false and adds a field error
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if it's not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinChars returns true if a value contains at least n characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches returns true if a value matches a provided compiled regular expression pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// PermittedValue returns true if a value is in a list of permitted values - using Generics
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
