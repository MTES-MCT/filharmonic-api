package cerbere

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/authentication"
)

type CerbereSso struct {
	config authentication.SsoConfig
}

func New(config authentication.SsoConfig) *CerbereSso {
	return &CerbereSso{
		config: config,
	}
}

type UserInfos struct {
	Email string `xml:"authenticationSuccess>user"`
}

func (sso *CerbereSso) ValidateTicket(ticket string) (string, error) {
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
	if userInfos.Email == "" {
		return "", authentication.ErrTicketValidationFailed
	}
	return userInfos.Email, nil
}
