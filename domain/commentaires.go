package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (s *Service) CreateCommentaire(ctx *UserContext, idInspection int64, commentaire models.Commentaire) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}
	ok, err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}

	id, err := s.repo.CreateCommentaire(ctx, idInspection, commentaire)
	if err != nil {
		return 0, err
	}
	// _, err = s.AddEvenement(ctx, models.Evenement{
	// 	Type:         models.CommentaireGeneral,
	// 	InspectionId: idInspection,
	// 	Data:         `{"commentaireId": ` + strconv.FormatInt(id, 10) + ` }`,
	// })
	// if err != nil {
	// 	return 0, err
	// }
	return id, err
}
