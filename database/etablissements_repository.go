package database

import (
	"github.com/MTES-MCT/filharmonic-api/database/helper"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	ErrEtablissementNotFound = errors.NewErrForbidden("Établissement non trouvé")
)

// utilisé seulement dans les tests
func (repo *Repository) ListEtablissements() ([]models.Etablissement, error) {
	etablissements := []models.Etablissement{}
	query := repo.db.client.Model(&etablissements)
	err := query.Select()
	return etablissements, err
}

func (repo *Repository) FindEtablissements(ctx *domain.UserContext, filter domain.ListEtablissementsFilter) (*models.FindEtablissementResults, error) {
	etablissements := []models.Etablissement{}

	query := repo.db.client.Model(&etablissements)
	if filter.S3IC != "" {
		query.Where("s3ic ilike ?", helper.BuildSearchValue(filter.S3IC))
	}
	if filter.Nom != "" {
		query.Where("f_unaccent(nom) ilike ? OR f_unaccent(raison) ilike ?", helper.BuildSearchValue(util.Normalize(filter.Nom)), helper.BuildSearchValue(util.Normalize(filter.Nom)))
	}
	if filter.Adresse != "" {
		adresse := helper.BuildSearchValue(util.Normalize(filter.Adresse))
		query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q.WhereOr("f_unaccent(adresse1) ilike ?", adresse).
				WhereOr("f_unaccent(adresse2) ilike ?", adresse).
				WhereOr("code_postal ilike ?", adresse).
				WhereOr("f_unaccent(commune) ilike ?", adresse)
			return q, nil
		})
	}

	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	total, err := query.Count()
	if err != nil {
		return nil, err
	}
	err = query.
		Limit(repo.config.PaginationSize).
		Offset((filter.GetPage() - 1) * repo.config.PaginationSize).
		Order("etablissement.nom").
		Select()
	return &models.FindEtablissementResults{
		Total:          total,
		Etablissements: etablissements,
	}, err
}

func (repo *Repository) GetEtablissementByID(ctx *domain.UserContext, id int64) (*models.Etablissement, error) {
	var etablissement models.Etablissement
	query := repo.db.client.Model(&etablissement).
		Relation("Inspections", func(q *orm.Query) (*orm.Query, error) {
			q.Order("id ASC")
			if ctx.IsExploitant() {
				q.ExcludeColumn("validation_rejetee", "motif_rejet_validation")
				q.Where("etat <> ?", models.EtatPreparation)
			}
			return q, nil
		}).
		Relation("Exploitants").
		Where("id = ?", id)
	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	}
	err := query.Select()
	if err == pg.ErrNoRows {
		return nil, ErrEtablissementNotFound
	}
	return &etablissement, err
}
