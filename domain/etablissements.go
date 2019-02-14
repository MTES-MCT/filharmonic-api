package domain

import "github.com/MTES-MCT/filharmonic-api/models"

type ListEtablissementsFilter struct {
	S3IC    string `form:"s3ic"`
	Nom     string `form:"nom"`
	Adresse string `form:"adresse"`
	Page    int    `form:"page"`
}

func (f *ListEtablissementsFilter) GetPage() int {
	if f.Page > 0 {
		return f.Page
	}
	return 1
}

func (s *Service) ListEtablissements(ctx *UserContext, filter ListEtablissementsFilter) (*models.FindEtablissementResults, error) {
	return s.repo.FindEtablissements(ctx, filter)
}

func (s *Service) GetEtablissement(ctx *UserContext, id int64) (*models.Etablissement, error) {
	return s.repo.GetEtablissementByID(ctx, id)
}
