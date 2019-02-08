package database

import (
	"strconv"
	"time"

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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.CreationPointDeControle,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"point_de_controle_id": ` + strconv.FormatInt(pointDeControle.Id, 10) + `}`,
		}
		err = tx.Insert(&evenement)
		if err != nil {
			return err
		}
		notification := models.Notification{
			EvenementId: evenement.Id,
		}
		err = tx.Insert(&notification)
		if err != nil {
			return err
		}
		return nil
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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.ModificationPointDeControle,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"point_de_controle_id": ` + strconv.FormatInt(pointDeControle.Id, 10) + `}`,
		}
		err = tx.Insert(&evenement)
		if err != nil {
			return err
		}
		notification := models.Notification{
			EvenementId: evenement.Id,
		}
		err = tx.Insert(&notification)
		if err != nil {
			return err
		}
		return nil
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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.PublicationPointDeControle,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"point_de_controle_id": ` + strconv.FormatInt(pointDeControle.Id, 10) + `}`,
		}
		err = tx.Insert(&evenement)
		if err != nil {
			return err
		}
		notification := models.Notification{
			EvenementId: evenement.Id,
		}
		err = tx.Insert(&notification)
		if err != nil {
			return err
		}
		return nil
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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.SuppressionPointDeControle,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"point_de_controle_id": ` + strconv.FormatInt(pointDeControle.Id, 10) + `}`,
		}
		err = tx.Insert(&evenement)
		if err != nil {
			return err
		}
		notification := models.Notification{
			EvenementId: evenement.Id,
		}
		err = tx.Insert(&notification)
		if err != nil {
			return err
		}
		return nil
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
func (repo *Repository) CheckEtatPointDeControle(id int64, etats []models.EtatInspection) (bool, error) {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspections AS i").
		JoinOn("i.id = point_de_controle.inspection_id").
		Where("i.etat in (?)", pg.In(etats)).
		Where("point_de_controle.id = ?", id).
		Count()
	return count == 1, err
}
