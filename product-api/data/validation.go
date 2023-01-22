package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice.
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new validation struct
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return &Validation{
		validate: validate,
	}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)

	if errs == nil {
		return nil
	}

	verrs := errs.(validator.ValidationErrors)

	var returnErrs ValidationErrors
	for _, err := range verrs {
		ve := ValidationError{err}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}

// validate SKU
func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}
