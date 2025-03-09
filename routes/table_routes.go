package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterTableRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	tableController := do.MustInvoke[*controllers.TableController](container)
	columnController := do.MustInvoke[*controllers.ColumnController](container)
	indexController := do.MustInvoke[*controllers.IndexController](container)

	tablesGroup := e.Group("api/tables", authMiddleware)

	// table routes
	tablesGroup.POST("", tableController.Store)
	tablesGroup.GET("", tableController.List)
	tablesGroup.GET("/:fullTableName", tableController.Show)
	tablesGroup.PUT("/:fullTableName/duplicate", tableController.Duplicate)
	tablesGroup.PUT("/:fullTableName/rename", tableController.Rename)
	tablesGroup.DELETE("/:fullTableName", tableController.Delete)

	// column routes
	tablesGroup.POST("/:fullTableName/columns", columnController.Store)
	tablesGroup.PUT("/:fullTableName/columns", columnController.Alter)
	tablesGroup.PUT("/:fullTableName/columns/:columnName", columnController.Rename)
	tablesGroup.DELETE("/:fullTableName/columns/:columnName", columnController.Delete)

	// index routes
	tablesGroup.POST("/:fullTableName/indexes", indexController.Store)
	tablesGroup.GET("/:fullTableName/indexes", indexController.List)
	tablesGroup.GET("/:fullTableName/indexes/:indexName", indexController.Show)
	tablesGroup.DELETE("/:fullTableName/indexes/:indexName", indexController.Delete)
}
