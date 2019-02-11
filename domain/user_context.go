package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

type UserContext struct {
	User *models.User
}

func (ctx *UserContext) IsInspecteur() bool {
	return ctx.User.Profile == models.ProfilInspecteur
}

func (ctx *UserContext) IsExploitant() bool {
	return ctx.User.Profile == models.ProfilExploitant
}

func (ctx *UserContext) IsApprobateur() bool {
	return ctx.User.Profile == models.ProfilApprobateur
}
