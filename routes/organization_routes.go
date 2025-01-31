package routes

import (
	"fluxton/controllers"
	"fluxton/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterOrganizationRoutes(e *echo.Echo, organizationController *controllers.OrganizationController) {
	organizationsGroup := e.Group("api/organizations", middleware.AuthMiddleware)

	organizationsGroup.POST("", organizationController.Store)
	organizationsGroup.GET("", organizationController.List)
	organizationsGroup.GET("/:id", organizationController.Show)
	organizationsGroup.PUT("/:id", organizationController.Update)
	organizationsGroup.DELETE("/:id", organizationController.Delete)
}
