package table_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fluxton/requests/column_requests"
	"fluxton/types"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type TableCreateRequest struct {
	requests.BaseRequest
	Name    string              `json:"name"`
	Columns []types.TableColumn `json:"columns"`
}

func (r *TableCreateRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var errors []string

	if err := r.validate(); err != nil {
		errors = append(errors, r.ExtractValidationErrors(err)...)

		return errors
	}

	for _, column := range r.Columns {
		if err := column_requests.ValidateColumn(column); err != nil {
			errors = append(errors, err.Error())

			return errors
		}
	}

	return errors
}

func (r *TableCreateRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.Name,
			validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscorePattern()),
			).Error("Table name must be alphanumeric with underscores"),
			validation.Length(
				configs.MinTableNameLength, configs.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"Name must be between %d and %d characters",
					configs.MinTableNameLength,
					configs.MaxTableNameLength,
				),
			),
			validation.By(validateName),
		),
		validation.Field(
			&r.Columns,
			validation.Required.Error("Columns are required"),
		),
	)
}
