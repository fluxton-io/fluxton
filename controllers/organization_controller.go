package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type OrganizationController struct {
	organizationService services.OrganizationService
}

func NewOrganizationController(injector *do.Injector) (*OrganizationController, error) {
	organizationService := do.MustInvoke[services.OrganizationService](injector)

	return &OrganizationController{organizationService: organizationService}, nil
}

func (nc *OrganizationController) List(c echo.Context) error {
	authUserId, _ := utils.NewAuth(c).Id()

	paginationParams := utils.ExtractPaginationParams(c)
	organizations, err := nc.organizationService.List(paginationParams, authUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResourceCollection(organizations))
}

func (nc *OrganizationController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organization, err := nc.organizationService.GetByID(id, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(&organization))
}

func (nc *OrganizationController) Store(c echo.Context) error {
	var request requests.OrganizationCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "organization.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	organization, err := nc.organizationService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.OrganizationResource(&organization))
}

func (nc *OrganizationController) Update(c echo.Context) error {
	var request requests.OrganizationCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "organization.error.invalidPayload")
	}

	updatedOrganization, err := nc.organizationService.Update(id, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(updatedOrganization))
}

func (nc *OrganizationController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := nc.organizationService.Delete(id, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
