package policies

import (
	"fluxton/models"
	"fluxton/repositories"
	"github.com/samber/do"
)

type ProjectPolicy struct {
	organizationRepo *repositories.OrganizationRepository
}

func NewProjectPolicy(injector *do.Injector) (*ProjectPolicy, error) {
	repo := do.MustInvoke[*repositories.OrganizationRepository](injector)

	return &ProjectPolicy{
		organizationRepo: repo,
	}, nil
}

func (s *ProjectPolicy) CanCreate(organizationID uint, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsLordOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanList(organizationID uint, authenticatedUserId uint) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUserId)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanView(organizationID uint, authenticatedUser models.AuthenticatedUser) bool {
	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}

func (s *ProjectPolicy) CanUpdate(organizationID uint, authenticatedUser models.AuthenticatedUser) bool {
	if !authenticatedUser.IsLordOrMore() {
		return false
	}

	isOrganizationUser, err := s.organizationRepo.IsOrganizationUser(organizationID, authenticatedUser.ID)
	if err != nil {
		return false
	}

	return isOrganizationUser
}
