package controllers

import (
	"fluxton/requests"
	"fluxton/requests/container_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type ContainerController struct {
	containerService services.ContainerService
}

func NewContainerController(injector *do.Injector) (*ContainerController, error) {
	containerService := do.MustInvoke[services.ContainerService](injector)

	return &ContainerController{containerService: containerService}, nil
}

// List retrieves all container
//
// @Summary List all container
// @Description Retrieve a list of container in a specified project.
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {object} responses.Response{content=[]resources.ContainerResponse} "List of container"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage [get]
func (bc *ContainerController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	container, err := bc.containerService.List(paginationParams, request.ProjectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ContainerResourceCollection(container))
}

// Show retrieves details of a specific container.
//
// @Summary Show details of a single container
// @Description Get details of a specific container
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
//
// @Success 200 {object} responses.Response{content=resources.ContainerResponse} "Container details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/containers/{containerUUID} [get]
func (bc *ContainerController) Show(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader

	authUser, _ := utils.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	container, err := bc.containerService.GetByUUID(containerUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ContainerResource(&container))
}

// Store creates a new container
//
// @Summary Create a new container
// @Description Add a new container to a project
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
// @Param container body container_requests.CreateRequest true "Container details"
//
// @Success 201 {object} responses.Response{content=resources.ContainerResponse} "Container created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage [post]
func (bc *ContainerController) Store(c echo.Context) error {
	var request container_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	container, err := bc.containerService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.ContainerResource(&container))
}

// Update a container
//
// @Summary Update a container
// @Description Modify an existing container's details
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param container body container_requests.CreateRequest true "Container details"
//
// @Success 200 {object} responses.Response{content=resources.ContainerResponse} "Container updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/containers/{containerUUID} [put]
func (bc *ContainerController) Update(c echo.Context) error {
	var request container_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedContainer, err := bc.containerService.Update(containerUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ContainerResource(updatedContainer))
}

// Delete a container
//
// @Summary Delete a container
// @Description Remove a container from a project
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
//
// @Success 204 "Container deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/containers/{containerUUID} [delete]
func (bc *ContainerController) Delete(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := bc.containerService.Delete(request, containerUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
