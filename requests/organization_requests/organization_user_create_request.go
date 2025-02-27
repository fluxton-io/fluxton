package organization_requests

import (
	"fluxton/requests"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationUserCreateRequest struct {
	requests.BaseRequest
	UserID uuid.UUID `json:"user_id"`
}

func (r *OrganizationUserCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required.Error("UserID is required")),
	)

	return r.ExtractValidationErrors(err)
}
