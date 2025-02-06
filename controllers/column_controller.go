package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type ColumnController struct {
	columnService services.ColumnService
}

func NewColumnController(injector *do.Injector) (*ColumnController, error) {
	columnService := do.MustInvoke[services.ColumnService](injector)

	return &ColumnController{columnService: columnService}, nil
}

func (pc *ColumnController) Store(c echo.Context) error {
	var request requests.ColumnCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, tableID, _, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := pc.columnService.Create(projectID, tableID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.TableResource(&table))
}

func (pc *ColumnController) Alter(c echo.Context) error {
	var request requests.ColumnAlterRequest
	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, columnName, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	renamedTable, err := pc.columnService.Alter(columnName, tableID, projectID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(renamedTable))
}

func (pc *ColumnController) Rename(c echo.Context) error {
	var request requests.ColumnRenameRequest
	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, columnName, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	renamedTable, err := pc.columnService.Rename(columnName, tableID, projectID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(renamedTable))
}

func (pc *ColumnController) Delete(c echo.Context) error {
	var request requests.DefaultRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, tableID, columnName, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.columnService.Delete(columnName, tableID, request.OrganizationID, projectID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (pc *ColumnController) parseRequest(c echo.Context) (uuid.UUID, uuid.UUID, string, error) {
	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, "", err
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, "", err
	}

	columnName := c.Param("columnName")

	return projectID, tableID, columnName, nil
}
