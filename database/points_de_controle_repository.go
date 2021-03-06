package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) CreatePointDeControle(ctx *domain.UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error) {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle.Id = 0
		pointDeControle.InspectionId = idInspection
		err := tx.Insert(&pointDeControle)
		if err != nil {
			return err
		}
		pointDeControle.Order = pointDeControle.Id
		_, err = tx.Model(&pointDeControle).Column("order").WherePK().Update()
		if err != nil {
			return err
		}

		err = repo.CreateEvenementTx(tx, ctx, models.EvenementCreationPointDeControle, idInspection, map[string]interface{}{
			"point_de_controle_id": pointDeControle.Id,
		})
		return err
	})
	return pointDeControle.Id, err
}

func (repo *Repository) UpdatePointDeControle(ctx *domain.UserContext, idPointDeControle int64, pointDeControle models.PointDeControle) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle.Id = idPointDeControle
		columns := []string{"sujet", "references_reglementaires"}
		_, err := tx.Model(&pointDeControle).Column(columns...).WherePK().Returning("inspection_id").Update()
		if err != nil {
			return err
		}

		err = repo.CreateEvenementTx(tx, ctx, models.EvenementModificationPointDeControle, pointDeControle.InspectionId, map[string]interface{}{
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return err
}

func (repo *Repository) PublishPointDeControle(ctx *domain.UserContext, idPointDeControle int64) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle := models.PointDeControle{
			Id:     idPointDeControle,
			Publie: true,
		}
		columns := []string{"publie"}
		_, err := tx.Model(&pointDeControle).Column(columns...).WherePK().Returning("inspection_id").Update()
		if err != nil {
			return err
		}

		err = repo.CreateEvenementTx(tx, ctx, models.EvenementPublicationPointDeControle, pointDeControle.InspectionId, map[string]interface{}{
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return err
}

func (repo *Repository) DeletePointDeControle(ctx *domain.UserContext, idPointDeControle int64) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle := models.PointDeControle{
			Id: idPointDeControle,
		}
		_, err := tx.Model(&pointDeControle).WherePK().Returning("inspection_id").Delete()
		if err != nil {
			return err
		}
		err = repo.CreateEvenementTx(tx, ctx, models.EvenementSuppressionPointDeControle, pointDeControle.InspectionId, map[string]interface{}{
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return err
}

func (repo *Repository) CheckUserAllowedPointDeControle(ctx *domain.UserContext, id int64) (bool, error) {
	if ctx.IsExploitant() {
		count, err := repo.db.client.Model(&models.PointDeControle{}).
			Join("JOIN inspections AS i").
			JoinOn("i.id = point_de_controle.inspection_id").
			Join("JOIN etablissements AS e").
			JoinOn("e.id = i.etablissement_id").
			Join("JOIN etablissement_to_exploitants AS ex").
			JoinOn("ex.etablissement_id = e.id").
			JoinOn("ex.user_id = ?", ctx.User.Id).
			Where("point_de_controle.id = ?", id).
			Count()
		return count == 1, err
	} else {
		count, err := repo.db.client.Model(&models.PointDeControle{}).
			Join("JOIN inspection_to_inspecteurs AS u").
			JoinOn("u.inspection_id = point_de_controle.inspection_id").
			JoinOn("u.user_id = ?", ctx.User.Id).
			Where("point_de_controle.id = ?", id).
			Count()
		return count == 1, err
	}
}

func (repo *Repository) GetEtatInspectionByPointDeControleID(idPointDeControle int64) (models.EtatInspection, error) {
	inspection := &models.Inspection{}
	err := repo.db.client.Model(&models.PointDeControle{}).
		Column("i.etat").
		Join("JOIN inspections AS i").
		JoinOn("i.id = point_de_controle.inspection_id").
		Where("point_de_controle.id = ?", idPointDeControle).
		Select(inspection)
	if err != nil {
		return models.EtatInconnu, err
	}
	return inspection.Etat, err
}

func (repo *Repository) CanCreatePointDeControle(ctx *domain.UserContext, idInspection int64) error {
	count, err := repo.db.client.Model(&models.Inspection{}).
		Where("id = ?", idInspection).
		Where("etat in (?)", pg.In([]models.EtatInspection{models.EtatPreparation, models.EtatEnCours})).
		Where("suite_id IS NULL").
		Count()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.ErrCreationPointDeControleImpossible
	}
	return nil
}

func (repo *Repository) CanUpdatePointDeControle(ctx *domain.UserContext, idPointDeControle int64) error {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspections as i ").
		JoinOn("i.id = point_de_controle.inspection_id").
		JoinOn("etat in (?)", pg.In([]models.EtatInspection{models.EtatPreparation, models.EtatEnCours})).
		JoinOn("suite_id IS NULL").
		Where("point_de_controle.id = ?", idPointDeControle).
		Count()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.ErrModificationPointDeControleImpossible
	}
	return nil
}

func (repo *Repository) OrderPointsDeControle(ctx *domain.UserContext, idInspection int64, pointsDeControleIds []int64) error {
	pointsDeControle := []interface{}{}
	for index, pointDeControleId := range pointsDeControleIds {
		pointsDeControle = append(pointsDeControle, &models.PointDeControle{
			Id:    pointDeControleId,
			Order: int64(index + 1),
		})
	}
	_, err := repo.db.client.Model(pointsDeControle...).Column("order").Update()
	return err
}
