package database

import (
	"strconv"
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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.CreationConstat,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"constat_id": ` + strconv.FormatInt(constat.Id, 10) + `, "point_de_controle_id": ` + strconv.FormatInt(idPointDeControle, 10) + `}`,
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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.SuppressionConstat,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"constat_id": ` + strconv.FormatInt(pointDeControle.ConstatId, 10) + `, "point_de_controle_id": ` + strconv.FormatInt(idPointDeControle, 10) + `}`,
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

func (repo *Repository) CheckCanCreateConstat(ctx *domain.UserContext, idPointDeControle int64) (bool, error) {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspections AS i").
		JoinOn("i.id = point_de_controle.inspection_id").
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = point_de_controle.inspection_id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("point_de_controle.id = ?", idPointDeControle).
		Where("point_de_controle.constat_id IS NULL").
		Where("i.etat = ?", models.EtatEnCours).
		Count()
	return count == 1, err
}

func (repo *Repository) CheckCanDeleteConstat(ctx *domain.UserContext, idPointDeControle int64) (bool, error) {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspections AS i").
		JoinOn("i.id = point_de_controle.inspection_id").
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = point_de_controle.inspection_id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("point_de_controle.id = ?", idPointDeControle).
		Where("point_de_controle.constat_id IS NOT NULL").
		Where("i.etat = ?", models.EtatEnCours).
		Count()
	return count == 1, err
}
