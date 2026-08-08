package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate runs struct-tag validation and returns a single human-readable
// error message combining all failed fields, or nil if the struct is valid.
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		messages = append(messages, fmt.Sprintf("%s failed '%s' validation", fieldErr.Field(), fieldErr.Tag()))
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}
