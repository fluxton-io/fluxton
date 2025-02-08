package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type OrganizationCreateRequest struct {
	Name string `json:"name"`
}

func (r *OrganizationCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(3, 100).Error("Title must be between 3 and 100 characters"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_ -]*$`)).Error("Name must start with a letter and contain only alphanumeric characters, underscores, spaces, and dashes"),
		),
	)

	// If no errors, return nil
	if err == nil {
		return nil
	}

	var errors []string
	if ve, ok := err.(validation.Errors); ok {
		for _, validationErr := range ve {
			errors = append(errors, validationErr.Error())
		}
	}

	return errors
}
