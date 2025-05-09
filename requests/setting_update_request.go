package requests

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type IndividualSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SettingUpdateRequest struct {
	BaseRequest
	Settings []IndividualSetting `json:"settings"`
}

func (r *SettingUpdateRequest) BindAndValidate(c echo.Context) []string {
	// Bind the JSON payload
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	r.SetContext(c)

	// Check if Settings slice is present
	if len(r.Settings) == 0 {
		return []string{"Settings required"}
	}

	var errors []string
	// Loop through each setting and validate individually
	for i, setting := range r.Settings {
		err := validation.ValidateStruct(&setting,
			validation.Field(&setting.Name, validation.Required.Error("Name is required")),
			validation.Field(&setting.Value, validation.Required.Error("Value is required")),
		)
		if err != nil {
			if ve, ok := err.(validation.Errors); ok {
				for field, validationErr := range ve {
					errors = append(errors,
						fmt.Sprintf("Setting[%d] - %s: %s", i, field, validationErr.Error()))
				}
			}
		}
	}

	return errors
}
