package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrCreationConstatImpossible          = errors.NewErrForbidden("Impossible de créer le constat")
	ErrSuppressionConstatImpossible       = errors.NewErrForbidden("Impossible de supprimer le constat")
	ErrBesoinEtatTraitementNonConformites = errors.NewErrForbidden("L'inspection doit être à l'étape de traitement des non-conformités")
	ErrBesoinTypeConstatNonConforme       = errors.NewErrForbidden("Impossible de résoudre un constat conforme")
)

func (s *Service) CreateConstat(ctx *UserContext, idPointDeControle int64, constat models.Constat) (int64, error) {
	if !ctx.IsInspecteur() {
		return 0, ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	err = s.repo.CanCreateConstat(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateConstat(ctx, idPointDeControle, constat)
}

func (s *Service) DeleteConstat(ctx *UserContext, idPointDeControle int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	err := s.repo.CanDeleteConstat(ctx, idPointDeControle)
	if err != nil {
		return err
	}

	return s.repo.DeleteConstat(ctx, idPointDeControle)
}

func (s *Service) ResolveConstat(ctx *UserContext, idPointDeControle int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	etatInspection, err := s.repo.GetEtatInspectionByPointDeControleID(idPointDeControle)
	if err != nil {
		return err
	}
	if etatInspection != models.EtatTraitementNonConformites {
		return ErrBesoinEtatTraitementNonConformites
	}
	typeConstat, err := s.repo.GetTypeConstatByPointDeControleID(idPointDeControle)
	if err != nil {
		return err
	}
	if typeConstat == models.TypeConstatConforme {
		return ErrBesoinTypeConstatNonConforme
	}

	return s.repo.ResolveConstat(ctx, idPointDeControle)
}
