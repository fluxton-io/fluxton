package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterProjectRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	projectController := do.MustInvoke[*controllers.ProjectController](container)

	projectsGroup := e.Group("api/projects", authMiddleware)

	projectsGroup.POST("", projectController.Store)
	projectsGroup.GET("", projectController.List)
	projectsGroup.GET("/:projectID", projectController.Show)
	projectsGroup.PUT("/:projectID", projectController.Update)
	projectsGroup.DELETE("/:projectID", projectController.Delete)
}
