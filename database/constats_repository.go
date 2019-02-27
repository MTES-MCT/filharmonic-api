package database

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) CreateConstat(ctx *domain.UserContext, idPointDeControle int64, constat models.Constat) (int64, error) {
	constatId := int64(0)
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		constat.Id = 0
		err := tx.Insert(&constat)
		if err != nil {
			return err
		}
		constatId = constat.Id
		pointDeControle := models.PointDeControle{
			Id:        idPointDeControle,
			ConstatId: constatId,
		}
		columns := []string{"constat_id"}
		_, err = tx.Model(&pointDeControle).Column(columns...).WherePK().Returning("inspection_id").Update()
		if err != nil {
			return err
		}
		err = repo.CreateEvenementTx(tx, ctx, models.EvenementCreationConstat, pointDeControle.InspectionId, map[string]interface{}{
			"constat_id":           constat.Id,
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return constatId, err
}

func (repo *Repository) DeleteConstat(ctx *domain.UserContext, idPointDeControle int64) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle := models.PointDeControle{
			Id: idPointDeControle,
		}
		err := tx.Model(&pointDeControle).Column("constat_id", "inspection_id").WherePK().Select()
		if err != nil {
			return err
		}
		constat := models.Constat{
			Id: pointDeControle.ConstatId,
		}
		err = tx.Delete(&constat)
		if err != nil {
			return err
		}
		err = repo.CreateEvenementTx(tx, ctx, models.EvenementSuppressionConstat, pointDeControle.InspectionId, map[string]interface{}{
			"constat_id":           pointDeControle.ConstatId,
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return err
}

func (repo *Repository) CanCreateConstat(ctx *domain.UserContext, idPointDeControle int64) error {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspections AS i").
		JoinOn("i.id = point_de_controle.inspection_id").
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = point_de_controle.inspection_id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("point_de_controle.id = ?", idPointDeControle).
		Where("point_de_controle.constat_id IS NULL").
		Where("i.etat = ?", models.EtatEnCours).
		Where("i.suite_id IS NULL").
		Count()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.ErrCreationConstatImpossible
	}
	return nil
}

func (repo *Repository) CanDeleteConstat(ctx *domain.UserContext, idPointDeControle int64) error {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspections AS i").
		JoinOn("i.id = point_de_controle.inspection_id").
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = point_de_controle.inspection_id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("point_de_controle.id = ?", idPointDeControle).
		Where("point_de_controle.constat_id IS NOT NULL").
		Where("i.etat = ?", models.EtatEnCours).
		Where("i.suite_id IS NULL").
		Count()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.ErrSuppressionConstatImpossible
	}
	return nil
}

func (repo *Repository) ResolveConstat(ctx *domain.UserContext, idPointDeControle int64) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle := models.PointDeControle{}
		_, err := tx.QueryOne(&pointDeControle, `UPDATE constats as c
											 SET date_resolution = ?
											 FROM
											 point_de_controles as p
											 WHERE
											 p.constat_id = c.id
											 AND p.id = ?
											 RETURNING p.inspection_id`, time.Now(), idPointDeControle)
		if err != nil {
			return err
		}

		err = repo.CreateEvenementTx(tx, ctx, models.EvenementResolutionConstat, pointDeControle.InspectionId, map[string]interface{}{
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return err
}

func (repo *Repository) GetTypeConstatByPointDeControleID(idPointDeControle int64) (models.TypeConstat, error) {
	constat := &models.Constat{}
	err := repo.db.client.Model(&models.PointDeControle{}).
		Column("c.type").
		Join("JOIN constats AS c").
		JoinOn("c.id = point_de_controle.constat_id").
		Where("point_de_controle.id = ?", idPointDeControle).
		Select(constat)
	if err != nil {
		return models.TypeConstatInconnu, err
	}
	return constat.Type, err
}
