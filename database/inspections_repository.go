package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) ListInspections(ctx *domain.UserContext) ([]models.Inspection, error) {
	var inspections []models.Inspection
	query := repo.db.client.Model(&inspections).Relation("Etablissement")
	if ctx.IsInspecteur() {
		query.Join("JOIN inspection_to_inspecteurs AS u").
			JoinOn("u.inspection_id = inspection.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	} else if ctx.IsExploitant() {
		query.Join("JOIN etablissements AS e").
			JoinOn("inspection.etablissement_id = e.id").
			Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = e.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	return inspections, err
}
