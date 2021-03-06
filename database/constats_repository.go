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
		columns := []string{"type", "remarques"}
		if constat.Type != models.TypeConstatConforme {
			columns = append(columns, "delai_nombre", "delai_unite")
		}
		_, err := tx.Model(&constat).Column(columns...).Returning("id").Insert()
		if err != nil {
			return err
		}
		constatId = constat.Id
		pointDeControle := models.PointDeControle{
			Id:        idPointDeControle,
			ConstatId: constatId,
		}
		columns = []string{"constat_id"}
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

func (repo *Repository) UpdateConstat(ctx *domain.UserContext, idPointDeControle int64, constat models.Constat) error {
	return repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		pointDeControle := models.PointDeControle{
			Id: idPointDeControle,
		}
		err := tx.Model(&pointDeControle).Column("constat_id", "inspection_id").WherePK().Select()
		if err != nil {
			return err
		}
		constat.Id = pointDeControle.ConstatId
		columns := []string{"type", "remarques", "delai_nombre", "delai_unite"}
		if constat.Type == models.TypeConstatConforme {
			constat.DelaiNombre = 0
			constat.DelaiUnite = ""
		}
		_, err = tx.Model(&constat).Column(columns...).WherePK().Update()
		if err != nil {
			return err
		}
		err = repo.CreateEvenementTx(tx, ctx, models.EvenementModificationConstat, pointDeControle.InspectionId, map[string]interface{}{
			"constat_id":           constat.Id,
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
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

func (repo *Repository) CanDeleteOrUpdateConstat(ctx *domain.UserContext, idPointDeControle int64) error {
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
		return domain.ErrSuppressionOuModificationConstatImpossible
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
