package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) CreateSuite(ctx *domain.UserContext, idInspection int64, suite models.Suite) (int64, error) {
	suiteId := int64(0)
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		suite.Id = 0
		err := tx.Insert(&suite)
		if err != nil {
			return err
		}
		suiteId = suite.Id
		inspection := models.Inspection{
			Id:      idInspection,
			SuiteId: suiteId,
		}
		columns := []string{"suite_id"}
		_, err = tx.Model(&inspection).Column(columns...).WherePK().Update()
		if err != nil {
			return err
		}
		err = repo.CreateEvenementTx(tx, ctx, models.EvenementCreationSuite, idInspection, nil)
		return err
	})
	return suiteId, err
}

func (repo *Repository) UpdateSuite(ctx *domain.UserContext, idInspection int64, suite models.Suite) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&suite).
			WherePK().Update()
		if err != nil {
			return err
		}

		err = repo.CreateEvenementTx(tx, ctx, models.EvenementModificationSuite, idInspection, nil)
		return err
	})
	return err
}

func (repo *Repository) DeleteSuite(ctx *domain.UserContext, idInspection int64) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		inspection := models.Inspection{
			Id: idInspection,
		}
		err := tx.Model(&inspection).Column("suite_id").WherePK().Select()
		if err != nil {
			return err
		}
		suite := models.Suite{
			Id: inspection.SuiteId,
		}
		err = tx.Delete(&suite)
		if err != nil {
			return err
		}

		err = repo.CreateEvenementTx(tx, ctx, models.EvenementSuppressionSuite, idInspection, nil)
		return err
	})
	return err
}

type InspectionWithNombrePointsDeControle struct {
	Count                         int
	NbPointsDeControleNonPublies  int
	NbPointsDeControleSansConstat int
}

func (repo *Repository) CheckCanCreateSuite(ctx *domain.UserContext, idInspection int64) error {
	results := &InspectionWithNombrePointsDeControle{}
	err := repo.db.client.Model(&models.Inspection{}).
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = inspection.id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("inspection.id = ?", idInspection).
		Where("inspection.suite_id IS NULL").
		Where("inspection.etat = ?", models.EtatEnCours).
		ColumnExpr("COUNT(inspection.id) as count").
		ColumnExpr(`(SELECT count(*)
			FROM point_de_controles AS point_de_controle
			WHERE point_de_controle.inspection_id = ?
			AND point_de_controle.deleted_at IS NULL
			AND point_de_controle.publie = FALSE) AS nb_points_de_controle_non_publies`, idInspection).
		ColumnExpr(`(SELECT count(*)
			FROM point_de_controles AS point_de_controle
			WHERE point_de_controle.inspection_id = ?
			AND point_de_controle.deleted_at IS NULL
			AND point_de_controle.constat_id IS NULL) AS nb_points_de_controle_sans_constat`, idInspection).
		Select(results)
	if err != nil {
		return err
	}
	if results.Count == 0 {
		return domain.ErrCreationSuiteImpossible
	}
	if results.NbPointsDeControleNonPublies > 0 {
		return models.ErrPointDeControleNonPublie
	}
	if results.NbPointsDeControleSansConstat > 0 {
		return models.ErrConstatManquant
	}
	return nil
}

func (repo *Repository) CheckCanDeleteSuite(ctx *domain.UserContext, idInspection int64) (bool, error) {
	count, err := repo.db.client.Model(&models.Inspection{}).
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = inspection.id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("inspection.id = ?", idInspection).
		Where("inspection.suite_id IS NOT NULL").
		Where("inspection.etat = ?", models.EtatEnCours).
		Count()
	return count == 1, err
}
