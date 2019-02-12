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
		err = repo.CreateEvenement(tx, ctx, models.EvenementCreationSuite, idInspection, nil)
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

		err = repo.CreateEvenement(tx, ctx, models.EvenementModificationSuite, idInspection, nil)
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

		err = repo.CreateEvenement(tx, ctx, models.EvenementSuppressionSuite, idInspection, nil)
		return err
	})
	return err
}

func (repo *Repository) CheckCanCreateSuite(ctx *domain.UserContext, idInspection int64) (bool, error) {
	count, err := repo.db.client.Model(&models.Inspection{}).
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = inspection.id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("inspection.id = ?", idInspection).
		Where("inspection.suite_id IS NULL").
		Where("inspection.etat = ?", models.EtatEnCours).
		Count()
	return count == 1, err
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
