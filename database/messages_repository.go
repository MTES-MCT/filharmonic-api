package database

import (
	"strconv"
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
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
			typeEvenement = models.CreationCommentaire
		} else {
			typeEvenement = models.CreationMessage
		}
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         typeEvenement,
			InspectionId: pointDeControle.InspectionId,
			Data:         `{"message_id": ` + strconv.FormatInt(messageId, 10) + `, "point_de_controle_id": ` + strconv.FormatInt(idPointDeControle, 10) + `}`,
		}
		err = tx.Insert(&evenement)
		if err != nil {
			return err
		}
		notification := models.Notification{
			EvenementId: evenement.Id,
		}
		err = tx.Insert(&notification)
		if err != nil {
			return err
		}
		return nil
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
		evenement := models.Evenement{
			AuteurId:     ctx.User.Id,
			CreatedAt:    time.Now(),
			Type:         models.LectureMessage,
			InspectionId: message.PointDeControle.InspectionId,
			Data:         `{"message_id": ` + strconv.FormatInt(idMessage, 10) + `, "point_de_controle_id": ` + strconv.FormatInt(message.PointDeControleId, 10) + `}`,
		}
		err = tx.Insert(&evenement)
		if err != nil {
			return err
		}
		notification := models.Notification{
			EvenementId: evenement.Id,
		}
		err = tx.Insert(&notification)
		if err != nil {
			return err
		}
		return nil
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
	profilAuteurs := make([]models.Profil, 0)
	if ctx.IsExploitant() {
		profilAuteurs = append(profilAuteurs, models.ProfilInspecteur, models.ProfilApprobateur)
	} else {
		profilAuteurs = append(profilAuteurs, models.ProfilExploitant)
	}

	count, err := repo.db.client.Model(&models.Message{}).
		Join("JOIN users AS u").
		JoinOn("u.id = message.auteur_id").
		Where("u.profile in (?)", pg.In(profilAuteurs)).
		Where("message.id = ?", id).
		Count()
	return count == 1, err
}
