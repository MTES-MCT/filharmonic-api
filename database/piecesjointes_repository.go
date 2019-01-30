package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func (repo *Repository) CreatePieceJointe(pieceJointe models.PieceJointe) (int64, error) {
	pieceJointe.Id = 0
	err := repo.db.client.Insert(&pieceJointe)
	return pieceJointe.Id, err
}

func (repo *Repository) GetPieceJointe(ctx *domain.UserContext, idPieceJointe int64) (*models.PieceJointe, error) {
	pieceJointe := models.PieceJointe{}
	var err error
	if ctx.IsExploitant() {
		err = repo.db.client.Model(&pieceJointe).Column("piece_jointe.*").
			Join("LEFT JOIN messages AS m").
			JoinOn("m.id = piece_jointe.message_id").
			Join("LEFT JOIN point_de_controles AS p").
			JoinOn("p.id = m.point_de_controle_id").
			Join("LEFT JOIN inspections AS i").
			JoinOn("i.id = p.inspection_id").
			Join("LEFT JOIN etablissements AS e").
			JoinOn("e.id = i.etablissement_id").
			Join("LEFT JOIN etablissement_to_exploitants AS ex").
			JoinOn("ex.etablissement_id = e.id").
			JoinOn("ex.user_id = ?", ctx.User.Id).
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.Where("(piece_jointe.message_id is NULL AND piece_jointe.commentaire_id is NULL AND piece_jointe.auteur_id = ?)", ctx.User.Id).
					WhereOr("(piece_jointe.message_id is not NULL AND piece_jointe.commentaire_id is not NULL AND piece_jointe.auteur_id <> ?)", ctx.User.Id)
				return q, nil
			}).
			Where("piece_jointe.id = ?", idPieceJointe).
			Select()
	} else {
		err = repo.db.client.Model(&pieceJointe).Column("piece_jointe.*").
			Join("LEFT JOIN commentaires AS c").
			JoinOn("c.id = piece_jointe.commentaire_id").
			Join("LEFT JOIN messages AS m").
			JoinOn("m.id = piece_jointe.message_id").
			Join("LEFT JOIN point_de_controles AS p").
			JoinOn("p.id = m.point_de_controle_id").
			Join("LEFT JOIN inspection_to_inspecteurs AS u").
			JoinOn("u.inspection_id = p.inspection_id").
			JoinOnOr("u.inspection_id = c.inspection_id").
			JoinOn("u.user_id = ?", ctx.User.Id).
			WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.Where("(piece_jointe.message_id is NULL AND piece_jointe.commentaire_id is NULL AND piece_jointe.auteur_id = ?)", ctx.User.Id).
					WhereOr("(piece_jointe.message_id is not NULL AND piece_jointe.commentaire_id is not NULL AND piece_jointe.auteur_id <> ?)", ctx.User.Id)
				return q, nil
			}).
			Where("piece_jointe.id = ?", idPieceJointe).
			Select()
	}
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &pieceJointe, err
}
