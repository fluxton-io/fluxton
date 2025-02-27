package organization_requests

import (
	"fluxton/requests"
	"fluxton/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type OrganizationCreateRequest struct {
	requests.BaseRequest
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
			validation.Length(3, 100).Error("Name must be between 3 and 100 characters"),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithSpaceUnderScoreAndDashPattern()),
			).Error("Organization name must be alphanumeric with underscores, spaces and dashes"),
		),
	)

	return r.ExtractValidationErrors(err)
}
