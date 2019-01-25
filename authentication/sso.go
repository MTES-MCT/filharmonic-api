package authentication

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
)

var (
	ErrUnauthorized = errors.New("unauthorized user")
	ErrMissingUser  = errors.New("missing user")
)

type Config struct {
	URL string `json:"url" default:"https://authentification.din.developpement-durable.gouv.fr/cas/public"`
}

type Sso struct {
	config Config
	repo   Repository
}

func New(config Config, repo Repository) *Sso {
	return &Sso{
		config: config,
		repo:   repo,
	}
}

func GenerateToken(id int64) string {
	return "token-" + strconv.FormatInt(id, 10)
}

type LoginResult struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (sso *Sso) Login(ticket string) (*LoginResult, error) {
	log.Debug().Str("ticket", ticket).Msg("login")

	email, err := sso.validateTicket(ticket)
	if err != nil {
		return nil, err
	}

	user, err := sso.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrMissingUser
	}
	result := &LoginResult{
		Token: sessions.New(user.Id),
		User:  *user,
	}
	return result, nil
}

func (sso *Sso) ValidateToken(token string) (*domain.UserContext, error) {
	userId := sessions.Get(token)
	if userId == 0 {
		return nil, ErrUnauthorized
	}
	user, err := sso.repo.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrMissingUser
	}
	return &domain.UserContext{User: user}, nil
}

type UserInfos struct {
	Email string `xml:"authenticationSuccess>user"`
}

func (sso *Sso) validateTicket(ticket string) (string, error) {
	res, err := http.Post(sso.config.URL+"/serviceValidate?ticket="+ticket, "text/xml", nil)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	userInfos := UserInfos{}
	err = xml.Unmarshal(data, &userInfos)
	if err != nil {
		return "", err
	}
	return userInfos.Email, nil
}
