package services

import (
	"github.com/samber/do"
	"myapp/errs"
	"myapp/models"
	"myapp/policies"
	"myapp/repositories"
	"myapp/requests"
	"myapp/utils"
)

type ProjectService interface {
	List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Project, error)
	GetByID(projectId uint, authenticatedUser models.AuthenticatedUser) (models.Project, error)
	Create(request *requests.ProjectCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Project, error)
	Update(projectId uint, authenticatedUser models.AuthenticatedUser, request *requests.ProjectCreateRequest) (*models.Project, error)
	Delete(projectId uint, authenticatedUser models.AuthenticatedUser) (bool, error)
}

type ProjectServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	projectRepo   *repositories.ProjectRepository
}

func NewProjectService(injector *do.Injector) (ProjectService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	repo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &ProjectServiceImpl{
		projectPolicy: policy,
		projectRepo:   repo,
	}, nil
}

func (s *ProjectServiceImpl) List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Project, error) {
	return s.projectRepo.ListForUser(paginationParams, authenticatedUserId)
}

func (s *ProjectServiceImpl) GetByID(projectId uint, authenticatedUser models.AuthenticatedUser) (models.Project, error) {
	if !s.projectPolicy.CanView(projectId, authenticatedUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return s.projectRepo.GetByIDForUser(projectId, authenticatedUser.ID)
}

func (s *ProjectServiceImpl) Create(request *requests.ProjectCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Project, error) {
	if !s.projectPolicy.CanCreate(authenticatedUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.createForbidden")
	}

	project := models.Project{
		Name: request.Name,
	}

	_, err := s.projectRepo.Create(&project, authenticatedUser.ID)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (s *ProjectServiceImpl) Update(projectId uint, authenticatedUser models.AuthenticatedUser, request *requests.ProjectCreateRequest) (*models.Project, error) {
	if !s.projectPolicy.CanUpdate(projectId, authenticatedUser) {
		return &models.Project{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	project, err := s.projectRepo.GetByIDForUser(projectId, authenticatedUser.ID)
	if err != nil {
		return nil, err
	}

	err = utils.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	return s.projectRepo.Update(projectId, &project)
}

func (s *ProjectServiceImpl) Delete(projectId uint, authenticatedUser models.AuthenticatedUser) (bool, error) {
	if !s.projectPolicy.CanUpdate(projectId, authenticatedUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	return s.projectRepo.Delete(projectId)
}
