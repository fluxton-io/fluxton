package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterTableRoutes(
	e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	TableController *controllers.TableController,
	ColumnController *controllers.ColumnController,
	IndexController *controllers.IndexController,
) {
	tablesGroup := e.Group("api/projects/:projectID/tables", authMiddleware)

	// table routes
	tablesGroup.POST("", TableController.Store)
	tablesGroup.GET("", TableController.List)
	tablesGroup.GET("/:tableID", TableController.Show)
	tablesGroup.PUT("/:tableID/duplicate", TableController.Duplicate)
	tablesGroup.PUT("/:tableID/rename", TableController.Rename)
	tablesGroup.DELETE("/:tableID", TableController.Delete)

	// column routes
	tablesGroup.POST("/:tableID/columns", ColumnController.Store)
	tablesGroup.PUT("/:tableID/columns", ColumnController.Alter)
	tablesGroup.PUT("/:tableID/columns/:columnName", ColumnController.Rename)
	tablesGroup.DELETE("/:tableID/columns/:columnName", ColumnController.Delete)

	// index routes
	tablesGroup.POST("/:tableID/indexes", IndexController.Store)
	tablesGroup.GET("/:tableID/indexes", IndexController.List)
	tablesGroup.GET("/:tableID/indexes/:indexName", IndexController.Show)
	tablesGroup.DELETE("/:tableID/indexes/:indexName", IndexController.Delete)
}
