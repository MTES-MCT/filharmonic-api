package database

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func (repo *Repository) CreateMessage(ctx *domain.UserContext, idPointDeControle int64, message models.Message) (int64, error) {
	messageId := int64(0)
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		message.Id = 0
		message.PointDeControleId = idPointDeControle
		message.AuteurId = ctx.User.Id
		message.Date = time.Now()
		message.Lu = false
		err := tx.Insert(&message)
		if err != nil {
			return err
		}
		messageId = message.Id
		for _, pieceJointe := range message.PiecesJointes {
			pieceJointe.MessageId = messageId
			ok, errCheck := repo.checkPieceJointeFree(tx, ctx, pieceJointe.Id)
			if errCheck != nil {
				return errCheck
			}
			if !ok {
				return domain.ErrInvalidInput
			}
			_, err = tx.Model(&pieceJointe).Column("message_id").WherePK().Update()
			if err != nil {
				return err
			}
		}
		pointDeControle := models.PointDeControle{
			Id: idPointDeControle,
		}
		err = tx.Model(&pointDeControle).WherePK().Select()
		if err != nil {
			return err
		}
		var typeEvenement models.TypeEvenement
		if message.Interne {
			typeEvenement = models.EvenementCreationCommentaire
		} else {
			typeEvenement = models.EvenementCreationMessage
		}
		err = repo.CreateEvenementTx(tx, ctx, typeEvenement, pointDeControle.InspectionId, map[string]interface{}{
			"message_id":           messageId,
			"point_de_controle_id": idPointDeControle,
		})
		return err
	})
	return messageId, err
}

func (repo *Repository) checkPieceJointeFree(tx *pg.Tx, ctx *domain.UserContext, idPieceJointe int64) (bool, error) {
	count, err := tx.Model(&models.PieceJointe{}).
		Where("auteur_id = ?", ctx.User.Id).
		Where("id = ?", idPieceJointe).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr("message_id is NULL").
				WhereOr("commentaire_id is NULL")
			return q, nil
		}).
		Count()
	return count == 1, err
}

func (repo *Repository) LireMessage(ctx *domain.UserContext, idMessage int64) error {
	err := repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		message := models.Message{
			Id: idMessage,
			Lu: true,
		}
		columns := []string{"lu"}
		_, err := tx.Model(&message).Column(columns...).WherePK().Update()
		if err != nil {
			return err
		}
		err = tx.Model(&message).Relation("PointDeControle").WherePK().Select()
		if err != nil {
			return err
		}
		err = repo.CreateEvenementTx(tx, ctx, models.EvenementLectureMessage, message.PointDeControle.InspectionId, map[string]interface{}{
			"message_id":           idMessage,
			"point_de_controle_id": message.PointDeControleId,
		})
		return err
	})
	return err
}

func (repo *Repository) CheckUserAllowedMessage(ctx *domain.UserContext, id int64) (bool, error) {
	if ctx.IsExploitant() {
		count, err := repo.db.client.Model(&models.Message{}).
			Join("JOIN point_de_controles AS p").
			JoinOn("p.id = message.point_de_controle_id").
			Join("JOIN inspections AS i").
			JoinOn("i.id = p.inspection_id").
			Join("JOIN etablissements AS e").
			JoinOn("e.id = i.etablissement_id").
			Join("JOIN etablissement_to_exploitants AS ex").
			JoinOn("ex.etablissement_id = e.id").
			JoinOn("ex.user_id = ?", ctx.User.Id).
			Where("message.id = ?", id).
			Count()
		return count == 1, err
	} else {
		count, err := repo.db.client.Model(&models.Message{}).
			Join("JOIN point_de_controles AS p").
			JoinOn("p.id = message.point_de_controle_id").
			Join("JOIN inspection_to_inspecteurs AS u").
			JoinOn("u.inspection_id = p.inspection_id").
			JoinOn("u.user_id = ?", ctx.User.Id).
			Where("message.id = ?", id).
			Count()
		return count == 1, err
	}
}

func (repo *Repository) CheckUserIsRecipient(ctx *domain.UserContext, id int64) (bool, error) {
	count, err := repo.db.client.Model(&models.Message{}).
		Join("JOIN users AS u").
		JoinOn("u.id = message.auteur_id").
		Where("u.profile in (?)", pg.In(getDestinataires(ctx))).
		Where("message.id = ?", id).
		Count()
	return count == 1, err
}

func getDestinataires(ctx *domain.UserContext) []models.Profil {
	profilDestinataires := make([]models.Profil, 0)
	if ctx.IsExploitant() {
		profilDestinataires = append(profilDestinataires, models.ProfilInspecteur, models.ProfilApprobateur)
	} else {
		profilDestinataires = append(profilDestinataires, models.ProfilExploitant)
	}
	return profilDestinataires
}

func (repo *Repository) ListNouveauxMessages() ([]domain.NouveauxMessagesUser, error) {
	nouveauxMessagesExploitants := []RowNouveauMessageUser{}
	err := repo.buildQueryNouveauxMessage(func(query *orm.Query) (*orm.Query, error) {
		return query.
			Join("JOIN etablissement_to_exploitants AS e2e").
			JoinOn(`e2e.user_id = "user".id`).
			Join("JOIN etablissements AS etablissement").
			JoinOn(`e2e.etablissement_id = etablissement.id`).
			Join("JOIN inspections AS inspection").
			JoinOn(`inspection.etablissement_id = etablissement.id`), nil
	}, []models.Profil{"inspecteur", "approbateur"}).
		Select(&nouveauxMessagesExploitants)
	if err != nil {
		return nil, err
	}

	nouveauxMessagesInspecteurs := []RowNouveauMessageUser{}
	err = repo.buildQueryNouveauxMessage(func(query *orm.Query) (*orm.Query, error) {
		return query.
			Join("JOIN inspection_to_inspecteurs AS i2i").
			JoinOn(`i2i.user_id = "user".id`).
			Join("JOIN inspections AS inspection").
			JoinOn(`i2i.inspection_id = inspection.id`).
			Join("JOIN etablissements AS etablissement").
			JoinOn(`inspection.etablissement_id = etablissement.id`), nil
	}, []models.Profil{"exploitant"}).
		Select(&nouveauxMessagesInspecteurs)
	if err != nil {
		return nil, err
	}

	return repo.toNouveauxMessageUsers(append(nouveauxMessagesExploitants, nouveauxMessagesInspecteurs...)), err
}

func (repo *Repository) toNouveauxMessageUsers(rows []RowNouveauMessageUser) []domain.NouveauxMessagesUser {
	nouveauxMessagesUsers := []domain.NouveauxMessagesUser{}
	if len(rows) == 0 {
		return nouveauxMessagesUsers
	}

	currentUser := domain.NouveauxMessagesUser{
		Destinataire: models.User{
			Email: rows[0].EmailDestinataire,
			Nom:   rows[0].NomDestinataire,
		},
		Messages: []domain.NouveauMessage{},
	}
	for _, row := range rows {
		if row.EmailDestinataire != currentUser.Destinataire.Email {
			nouveauxMessagesUsers = append(nouveauxMessagesUsers, currentUser)
			currentUser = domain.NouveauxMessagesUser{
				Destinataire: models.User{
					Email: row.EmailDestinataire,
					Nom:   row.NomDestinataire,
				},
				Messages: []domain.NouveauMessage{},
			}
		}
		currentUser.Messages = append(currentUser.Messages, domain.NouveauMessage{
			DateInspection:       row.DateInspection,
			RaisonEtablissement:  row.RaisonEtablissement,
			SujetPointDeControle: row.SujetPointDeControle,
			Message:              row.Message,
			DateMessage:          util.FormatDateTime(row.DateMessage),
			NomAuteur:            row.NomAuteur,
			InspectionId:         row.InspectionId,
			PointDeControleId:    row.PointDeControleId,
			MessageId:            row.MessageId,
		})
	}
	nouveauxMessagesUsers = append(nouveauxMessagesUsers, currentUser)

	return nouveauxMessagesUsers
}

type RowNouveauMessageUser struct {
	EmailDestinataire    string
	NomDestinataire      string
	InspectionId         int64
	DateInspection       string
	RaisonEtablissement  string
	PointDeControleId    int64
	SujetPointDeControle string
	Message              string
	MessageId            int64
	DateMessage          time.Time
	NomAuteur            string
}

func (repo *Repository) buildQueryNouveauxMessage(joinFunc func(*orm.Query) (*orm.Query, error), profils []models.Profil) *orm.Query {
	return repo.db.client.Model(&models.User{}).
		ColumnExpr(`"user".email as email_destinataire`).
		ColumnExpr(`"user".prenom || ' ' || "user".nom as nom_destinataire`).
		ColumnExpr("inspection.date as date_inspection").
		ColumnExpr(`inspection.id as inspection_id`).
		ColumnExpr(`etablissement.raison as raison_etablissement`).
		ColumnExpr(`p.id as point_de_controle_id`).
		ColumnExpr(`p.sujet as sujet_point_de_controle`).
		ColumnExpr(`m.message as message`).
		ColumnExpr(`m.id as message_id`).
		ColumnExpr(`m.date as date_message`).
		ColumnExpr(`auteur.prenom || ' ' || auteur.nom as nom_auteur`).
		Apply(joinFunc).
		Join("JOIN point_de_controles AS p").
		JoinOn("p.inspection_id = inspection.id").
		JoinOn("p.publie IS TRUE").
		Join("JOIN messages AS m").
		JoinOn("m.point_de_controle_id = p.id").
		JoinOn("m.lu IS FALSE").
		JoinOn("m.interne IS FALSE").
		Join("JOIN users AS auteur").
		JoinOn("auteur.id = m.auteur_id").
		JoinOn("auteur.profile in (?)", pg.In(profils)).
		Order("user.id")
}
