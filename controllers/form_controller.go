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

func (pc *FormController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, "Invalid project UUID")
	}

	paginationParams := utils.ExtractPaginationParams(c)
	forms, err := pc.formService.List(paginationParams, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResourceCollection(forms))
}

func (pc *FormController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	form, err := pc.formService.GetByUUID(formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResource(&form))
}

func (pc *FormController) Store(c echo.Context) error {
	var request form_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	form, err := pc.formService.Create(projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormResource(&form))
}

func (pc *FormController) Update(c echo.Context) error {
	var request form_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := pc.formService.Update(formUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResource(updatedOrganization))
}

func (pc *FormController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.formService.Delete(formUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
