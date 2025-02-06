package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type ProjectUpdateRequest struct {
	Name string `json:"name"`
}

func (r *ProjectUpdateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Name is required"), validation.Length(3, 100).Error("Name must be between 3 and 100 characters")),
	)

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
