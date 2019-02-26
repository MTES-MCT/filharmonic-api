package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrPieceJointeNotFound = errors.NewErrForbidden("Pièce-jointe non trouvée")
)

func (s *Service) CreatePieceJointe(ctx *UserContext, pieceJointeFile models.PieceJointeFile) (int64, error) {
	storageId, err := s.storage.Put(pieceJointeFile)
	if err != nil {
		return 0, err
	}
	pieceJointe := models.PieceJointe{
		Nom:       pieceJointeFile.Nom,
		Type:      pieceJointeFile.Type,
		Taille:    pieceJointeFile.Taille,
		StorageId: storageId,
		AuteurId:  ctx.User.Id,
	}
	pieceJointeId, err := s.repo.CreatePieceJointe(pieceJointe)
	if err != nil {
		return 0, err
	}
	return pieceJointeId, nil
}

func (s *Service) GetPieceJointe(ctx *UserContext, idPieceJointe int64) (*models.PieceJointeFile, error) {
	pieceJointe, err := s.repo.GetPieceJointe(ctx, idPieceJointe)
	if err != nil {
		return nil, err
	}
	if pieceJointe == nil {
		return nil, ErrPieceJointeNotFound
	}
	reader, err := s.storage.Get(pieceJointe.StorageId)
	if err != nil {
		return nil, err
	}
	return &models.PieceJointeFile{
		Nom:     pieceJointe.Nom,
		Type:    pieceJointe.Type,
		Taille:  pieceJointe.Taille,
		Content: reader,
	}, nil
}
