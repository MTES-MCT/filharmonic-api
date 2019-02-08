package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) ListEvenements(ctx *domain.UserContext, filter domain.ListEvenementsFilter) ([]models.Evenement, error) {
	evenements := []models.Evenement{}
	err := repo.db.client.Model(&evenements).
		Relation("Auteur").
		Where("auteur_id = ?", ctx.User.Id).
		Select()
	return evenements, err
}

func (repo *Repository) GetEvenementByID(ctx *domain.UserContext, id int64) (*models.Evenement, error) {
	var evenement models.Evenement
	evenement.Id = id
	err := repo.db.client.Model(&evenement).
		Relation("Auteur").
		WherePK().
		Where("auteur_id = ?", ctx.User.Id).
		Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &evenement, err
}
