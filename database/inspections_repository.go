package database

import (
	"errors"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func (repo *Repository) ListInspections(ctx *domain.UserContext, filter domain.ListInspectionsFilter) ([]models.Inspection, error) {
	inspections := []models.Inspection{}
	query := repo.db.client.Model(&inspections).Relation("Etablissement")
	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	} else {
		if filter.Assigned {
			query.Join("JOIN inspection_to_inspecteurs AS u").
				JoinOn("u.inspection_id = inspection.id").
				JoinOn("u.user_id = ?", ctx.User.Id)
		}
	}
	err := query.Select()
	return inspections, err
}

func (repo *Repository) CreateInspection(ctx *domain.UserContext, inspection models.Inspection) (int64, error) {
	inspectionId := int64(0)
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		inspection.Id = 0
		inspection.Etat = models.EtatPreparation
		err := tx.Insert(&inspection)
		if err != nil {
			return err
		}

		for _, inspecteur := range inspection.Inspecteurs {
			err = tx.Insert(&models.InspectionToInspecteur{
				InspectionId: inspection.Id,
				UserId:       inspecteur.Id,
			})
			if err != nil {
				return err
			}
		}
		inspectionId = inspection.Id
		err = repo.CreateEvenement(tx, ctx, models.EvenementCreationInspection, inspectionId, nil)
		return err
	})
	return inspectionId, err
}

func (repo *Repository) UpdateInspection(ctx *domain.UserContext, inspection models.Inspection) error {
	return repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		columns := []string{"date", "type", "origine", "annonce", "circonstance", "detail_circonstance", "contexte", "themes"}
		_, err := tx.Model(&inspection).Column(columns...).WherePK().Update()
		if err != nil {
			return err
		}
		_, err = tx.Model(&models.InspectionToInspecteur{}).Where("inspection_id = ?", inspection.Id).Delete()
		if err != nil {
			return err
		}
		for _, inspecteur := range inspection.Inspecteurs {
			err = tx.Insert(&models.InspectionToInspecteur{
				InspectionId: inspection.Id,
				UserId:       inspecteur.Id,
			})
			if err != nil {
				return err
			}
		}
		err = repo.CreateEvenement(tx, ctx, models.EvenementModificationInspection, inspection.Id, nil)
		return err
	})
}

func (repo *Repository) GetInspectionByID(ctx *domain.UserContext, id int64) (*models.Inspection, error) {
	var inspection models.Inspection
	query := repo.db.client.Model(&inspection).
		Relation("Etablissement").
		Relation("PointsDeControle", func(q *orm.Query) (*orm.Query, error) {
			if ctx.IsExploitant() {
				q.Where("publie = TRUE")
			}
			return q, nil
		}).
		Relation("PointsDeControle.Constat").
		Relation("PointsDeControle.Messages", func(q *orm.Query) (*orm.Query, error) {
			if ctx.IsExploitant() {
				q.Where("interne = FALSE")
			}
			return q.Order("date ASC"), nil
		}).
		Relation("PointsDeControle.Messages.Auteur").
		Relation("PointsDeControle.Messages.PiecesJointes").
		Relation("Inspecteurs").
		Relation("Suite").
		Where(`inspection.id = ?`, id)
	if ctx.IsExploitant() {
		query.Join("JOIN etablissement_to_exploitants AS u").
			JoinOn("u.etablissement_id = etablissement.id").
			JoinOn("u.user_id = ?", ctx.User.Id)
	} else {
		query.Relation("Commentaires", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("date ASC"), nil
		}).
			Relation("Commentaires.Auteur").
			Relation("Commentaires.PiecesJointes").
			Relation("Evenements").
			Relation("Evenements.Auteur")
	}
	err := query.Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &inspection, err
}

func (repo *Repository) CheckInspecteurAllowedInspection(ctx *domain.UserContext, id int64) (bool, error) {
	count, err := repo.db.client.Model(&models.InspectionToInspecteur{}).
		Where("inspection_id = ?", id).
		Where("user_id = ?", ctx.User.Id).
		Count()
	return count == 1, err
}

func (repo *Repository) CheckEtatInspection(id int64, etats []models.EtatInspection) (bool, error) {
	count, err := repo.db.client.Model(&models.Inspection{}).
		Where("id = ?", id).
		Where("etat in (?)", pg.In(etats)).
		Count()
	return count == 1, err
}

func (repo *Repository) UpdateEtatInspection(ctx *domain.UserContext, id int64, etat models.EtatInspection) error {
	inspection := models.Inspection{
		Id:   id,
		Etat: etat,
	}
	columns := []string{"etat"}

	return repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&inspection).Column(columns...).WherePK().Update()
		if err != nil {
			return err
		}
		var typeEvenement models.TypeEvenement
		switch etat {
		case models.EtatEnCours:
			typeEvenement = models.EvenementPublicationInspection
		case models.EtatAttenteValidation:
			typeEvenement = models.EvenementDemandeValidationInspection
		case models.EtatValide:
			typeEvenement = models.EvenementValidationInspection
		case models.EtatNonValide:
			typeEvenement = models.EvenementRejetValidationInspection
		default:
			err = errors.New("etat unknown")
			return err
		}
		err = repo.CreateEvenement(tx, ctx, typeEvenement, inspection.Id, nil)
		return err
	})
}

func (repo *Repository) AddFavoriToInspection(ctx *domain.UserContext, idInspection int64) error {
	return repo.db.client.Insert(&models.UserToFavori{
		InspectionId: idInspection,
		UserId:       ctx.User.Id,
	})
}

func (repo *Repository) RemoveFavoriToInspection(ctx *domain.UserContext, idInspection int64) error {
	return repo.db.client.Delete(&models.UserToFavori{
		InspectionId: idInspection,
		UserId:       ctx.User.Id,
	})
}
