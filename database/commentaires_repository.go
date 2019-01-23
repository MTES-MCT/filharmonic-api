package database

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) CreateCommentaire(ctx *domain.UserContext, idInspection int64, commentaire models.Commentaire) (int64, error) {
	commentaire.Id = 0
	commentaire.InspectionId = idInspection
	commentaire.AuteurId = ctx.User.Id
	commentaire.Date = time.Now()
	err := repo.db.client.Insert(&commentaire)
	if err != nil {
		return 0, err
	}
	return commentaire.Id, nil
}
