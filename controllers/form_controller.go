package controllers

import (
	"fluxton/requests/form_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormController struct {
	formService services.FormService
}

func NewFormController(injector *do.Injector) (*FormController, error) {
	formService := do.MustInvoke[services.FormService](injector)

	return &FormController{formService: formService}, nil
}

// List retrieves all forms for a project
//
// @Summary List all forms
// @Description Retrieve a list of all forms for the specified project
// @Tags forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectId path string true "Project ID"
//
// @Success 200 {array} responses.Response{content=[]resources.FormResponse} "List of forms"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectId}/forms [get]
func (fc *FormController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, "Invalid project UUID")
	}

	paginationParams := utils.ExtractPaginationParams(c)
	forms, err := fc.formService.List(paginationParams, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResourceCollection(forms))
}

func (fc *FormController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	form, err := fc.formService.GetByUUID(formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResource(&form))
}

// Store creates a new form
//
// @Summary Create a new form
// @Description Add a new form with a name and description
// @Tags forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param form body form_requests.CreateRequest true "Form name and description"
//
// @Success 201 {object} responses.Response{content=resources.FormResponse} "Form created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms [post]
func (fc *FormController) Store(c echo.Context) error {
	var request form_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	form, err := fc.formService.Create(projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormResource(&form))
}

// Update updates an existing form
//
// @Summary Update an existing form
// @Description Update form details such as name and description
// @Tags forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param form body form_requests.CreateRequest true "Form name and description"
// @Param formId path string true "Form ID"
//
// @Success 200 {object} responses.Response{content=resources.FormResponse} "Form updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formId} [put]
func (fc *FormController) Update(c echo.Context) error {
	var request form_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedForm, err := fc.formService.Update(formUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResource(updatedForm))
}

func (fc *FormController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := fc.formService.Delete(formUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
