package di

import (
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"myapp/controllers"
	"myapp/db"
	"myapp/factories"
	"myapp/policies"
	"myapp/repositories"
	"myapp/services"
)

func InitializeContainer() *do.Injector {
	injector := do.New()

	// Database
	db.InitDB()
	do.Provide(injector, func(i *do.Injector) (*sqlx.DB, error) {
		return db.GetDB(), nil
	})

	// Repositories
	do.Provide(injector, repositories.NewUserRepository)
	do.Provide(injector, repositories.NewOrganizationRepository)
	do.Provide(injector, repositories.NewProjectRepository)
	do.Provide(injector, repositories.NewNoteRepository)
	do.Provide(injector, repositories.NewTagRepository)

	// Factories
	do.Provide(injector, factories.NewUserFactory)
	do.Provide(injector, factories.NewNoteFactory)
	do.Provide(injector, factories.NewTagFactory)

	// policies
	do.Provide(injector, policies.NewOrganizationPolicy)
	do.Provide(injector, policies.NewProjectPolicy)

	// Services
	do.Provide(injector, services.NewUserService)
	do.Provide(injector, services.NewNoteService)
	do.Provide(injector, services.NewOrganizationService)
	do.Provide(injector, services.NewProjectService)

	// Controllers
	do.Provide(injector, controllers.NewUserController)
	do.Provide(injector, controllers.NewNoteController)
	do.Provide(injector, controllers.NewOrganizationController)
	do.Provide(injector, controllers.NewProjectController)

	return injector
}
