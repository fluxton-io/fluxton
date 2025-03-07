package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/project_requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"math/rand"
	"strings"
	"time"
)

type ProjectService interface {
	List(paginationParams utils.PaginationParams, organizationUUID uuid.UUID, authUser models.AuthUser) ([]models.Project, error)
	GetByID(projectUUID uuid.UUID, authUser models.AuthUser) (models.Project, error)
	Create(request *project_requests.CreateRequest, authUser models.AuthUser) (models.Project, error)
	Update(projectUUID uuid.UUID, authUser models.AuthUser, request *project_requests.UpdateRequest) (*models.Project, error)
	Delete(projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type ProjectServiceImpl struct {
	projectPolicy    *policies.ProjectPolicy
	databaseRepo     *repositories.DatabaseRepository
	projectRepo      *repositories.ProjectRepository
	postgrestService PostgrestService
}

func NewProjectService(injector *do.Injector) (ProjectService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	postgrestService, _ := NewPostgrestService()

	return &ProjectServiceImpl{
		projectPolicy:    policy,
		databaseRepo:     databaseRepo,
		projectRepo:      projectRepo,
		postgrestService: postgrestService,
	}, nil
}

func (s *ProjectServiceImpl) List(paginationParams utils.PaginationParams, organizationUUID uuid.UUID, authUser models.AuthUser) ([]models.Project, error) {
	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Project{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authUser.Uuid)
}

func (s *ProjectServiceImpl) GetByID(projectUUID uuid.UUID, authUser models.AuthUser) (models.Project, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return models.Project{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return project, nil
}

func (s *ProjectServiceImpl) Create(request *project_requests.CreateRequest, authUser models.AuthUser) (models.Project, error) {
	if !s.projectPolicy.CanCreate(request.OrganizationUUID, authUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.createForbidden")
	}

	err := s.validateNameForDuplication(request.Name, request.OrganizationUUID)
	if err != nil {
		return models.Project{}, err
	}

	project := models.Project{
		Name:             request.Name,
		OrganizationUuid: request.OrganizationUUID,
		DBName:           s.generateDBName(),
		DBPort:           s.generateDBPort(),
		CreatedBy:        authUser.Uuid,
		UpdatedBy:        authUser.Uuid,
	}

	_, err = s.projectRepo.Create(&project)
	if err != nil {
		return models.Project{}, err
	}

	err = s.databaseRepo.Create(project.DBName)
	if err != nil {
		// TODO: handle better
		s.projectRepo.Delete(project.Uuid)

		return models.Project{}, err
	}

	err = s.postgrestService.StartContainer(project.DBName, project.DBPort)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (s *ProjectServiceImpl) Update(projectUUID uuid.UUID, authUser models.AuthUser, request *project_requests.UpdateRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &models.Project{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	err = utils.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	project.UpdatedAt = time.Now()
	project.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, project.OrganizationUuid)
	if err != nil {
		return &models.Project{}, err
	}

	return s.projectRepo.Update(&project)
}

func (s *ProjectServiceImpl) Delete(projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	err = s.databaseRepo.DropIfExists(project.DBName)
	if err != nil {
		return false, err
	}

	err = s.postgrestService.RemoveContainer(project.DBName)
	if err != nil {
		return false, err
	}

	return s.projectRepo.Delete(projectUUID)
}

func (s *ProjectServiceImpl) generateDBName() string {
	return "udb_" + strings.ReplaceAll(strings.ToLower(uuid.New().String()), "-", "")
}

func (s *ProjectServiceImpl) generateDBPort() int {
	return rand.Intn(65535-5000+1) + 5000
}

func (s *ProjectServiceImpl) validateNameForDuplication(name string, organizationUUID uuid.UUID) error {
	exists, err := s.projectRepo.ExistsByNameForOrganization(name, organizationUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("project.error.duplicateName")
	}

	return nil
}
