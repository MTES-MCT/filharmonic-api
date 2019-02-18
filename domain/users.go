package domain

import "github.com/MTES-MCT/filharmonic-api/models"

type ListUsersFilters struct {
	Inspecteurs bool `form:"inspecteurs"`
	Approbateurs bool `form:"approbateurs"`
}

func (s *Service) ListUsers(ctx *UserContext, filters ListUsersFilters) ([]models.User, error) {
	if ctx.IsExploitant() {
		return nil, ErrBesoinProfilInspecteur
	}
	return s.repo.FindUsers(filters)
}
