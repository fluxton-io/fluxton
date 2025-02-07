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

type OrganizationUserController struct {
	organizationService services.OrganizationService
}

func NewOrganizationUserController(injector *do.Injector) (*OrganizationUserController, error) {
	organizationService := do.MustInvoke[services.OrganizationService](injector)

	return &OrganizationUserController{organizationService: organizationService}, nil
}

func (nc *OrganizationUserController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationID, err := utils.GetUUIDPathParam(c, "organizationID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organizationUsers, err := nc.organizationService.ListUsers(organizationID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.UserResourceCollection(organizationUsers))
}

func (nc *OrganizationUserController) Store(c echo.Context) error {
	var request requests.OrganizationUserCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "organization.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	organizationID, err := utils.GetUUIDPathParam(c, "organizationID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organizationUser, err := nc.organizationService.CreateUser(&request, organizationID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.UserResource(&organizationUser))
}

func (nc *OrganizationUserController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationID, err := utils.GetUUIDPathParam(c, "organizationID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	userID, err := utils.GetUUIDPathParam(c, "userID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := nc.organizationService.DeleteUser(organizationID, userID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
