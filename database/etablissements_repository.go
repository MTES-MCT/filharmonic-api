package database

import (
	"strings"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func (repo *Repository) ListEtablissements() ([]models.Etablissement, error) {
	etablissements := []models.Etablissement{}
	query := repo.db.client.Model(&etablissements)
	err := query.Select()
	return etablissements, err
}

func (repo *Repository) FindEtablissements(ctx *domain.UserContext, filter domain.ListEtablissementsFilter) ([]models.Etablissement, error) {
	etablissements := []models.Etablissement{}

	query := repo.db.client.Model(&etablissements)
	if filter.S3IC != "" {
		query.Where("lower(s3ic) like ?", "%"+strings.ToLower(filter.S3IC)+"%")
	}
	if filter.Nom != "" {
		query.Where("lower(nom) like ? OR lower(raison) like ?", "%"+strings.ToLower(filter.Nom)+"%", "%"+strings.ToLower(filter.Nom)+"%")
	}
	if filter.Adresse != "" {
		adresseLowerCase := strings.ToLower(filter.Adresse)
		query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q.WhereOr("lower(adresse1) like ?", "%"+adresseLowerCase+"%").
				WhereOr("lower(adresse2) like ?", "%"+adresseLowerCase+"%").
				WhereOr("lower(code_postal) like ?", "%"+adresseLowerCase+"%").
				WhereOr("lower(commune) like ?", "%"+adresseLowerCase+"%")
			return q, nil
		})
	}

	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	return etablissements, err
}

func (repo *Repository) GetEtablissementByID(ctx *domain.UserContext, id int64) (*models.Etablissement, error) {
	var etablissement models.Etablissement
	query := repo.db.client.Model(&etablissement).
		Relation("Inspections").
		Relation("Exploitants").
		Where("id = ?", id)
	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &etablissement, err
}
