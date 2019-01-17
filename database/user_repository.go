package database

import (
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) GetUser(email string) (*models.User, error) {
	var user models.User
	err := repo.db.client.Model(&user).Where("email = ?", email).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &user, err
}
