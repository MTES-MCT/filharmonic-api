package database

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) CreateCommentaire(ctx *domain.UserContext, idInspection int64, commentaire models.Commentaire) (int64, error) {
	commentaireId := int64(0)
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		commentaire.Id = 0
		commentaire.InspectionId = idInspection
		commentaire.AuteurId = ctx.User.Id
		commentaire.Date = time.Now()
		err := tx.Insert(&commentaire)
		if err != nil {
			return err
		}
		commentaireId = commentaire.Id
		for _, pieceJointe := range commentaire.PiecesJointes {
			pieceJointe.CommentaireId = commentaireId
			ok, errCheck := repo.checkPieceJointeFree(tx, ctx, pieceJointe.Id)
			if errCheck != nil {
				return errCheck
			}
			if !ok {
				return domain.ErrInvalidInput
			}
			_, err = tx.Model(&pieceJointe).Column("commentaire_id").WherePK().Update()
			if err != nil {
				return err
			}
		}
		err = repo.CreateEvenement(tx, ctx, models.EvenementCommentaireGeneral, idInspection, map[string]interface{}{
			"commentaire_id":           commentaireId,
		})
		return err
	})
	return commentaireId, err
}
