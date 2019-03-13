package tests

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication"
	authmocks "github.com/MTES-MCT/filharmonic-api/authentication/mocks"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/domain"
	domainmocks "github.com/MTES-MCT/filharmonic-api/domain/mocks"
	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/MTES-MCT/filharmonic-api/events"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/MTES-MCT/filharmonic-api/templates"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func Init(t *testing.T) (*httpexpect.Expect, func()) {
	e, close, _ := InitWithSso(t)
	return e, close
}

func InitService(t *testing.T) (*require.Assertions, *app.Application, func()) {
	assert, a := InitDB(t)

	var err error
	a.Storage, err = storage.New(a.Config.Storage)
	assert.NoError(err)
	a.Sessions = sessions.NewMemory()
	sso := new(authmocks.Sso)
	a.Sso = sso
	a.AuthenticationService = authentication.New(a.Repo, a.Sso, a.Sessions)
	a.TemplateService, err = templates.New(a.Config.Templates)
	assert.NoError(err)
	emailService := new(domainmocks.EmailService)
	emailService.On("Send", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		email := args.Get(0).(emails.Email)
		assert.NoError(ioutil.WriteFile("../../.tmp/email-"+strconv.FormatInt(time.Now().UnixNano(), 10)+".html", []byte(email.HTMLPart), 0644))
	})
	a.EmailService = emailService
	a.Service = domain.New(a.Config.Service, a.Repo, a.Storage, a.TemplateService, a.EmailService)

	return assert, a, func() {
		err := a.Shutdown()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
}

func InitWithSso(t *testing.T) (*httpexpect.Expect, func(), *authmocks.Sso) {
	assert, a := InitDB(t)

	var err error
	a.Storage, err = storage.New(a.Config.Storage)
	assert.NoError(err)
	a.Sessions = sessions.NewMemory()
	sso := new(authmocks.Sso)
	a.Sso = sso
	a.AuthenticationService = authentication.New(a.Repo, a.Sso, a.Sessions)
	a.TemplateService, err = templates.New(a.Config.Templates)
	assert.NoError(err)
	emailService := new(domainmocks.EmailService)
	emailService.On("Send", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		email := args.Get(0).(emails.Email)
		assert.NoError(ioutil.WriteFile("../../.tmp/email-"+strconv.FormatInt(time.Now().UnixNano(), 10)+".html", []byte(email.HTMLPart), 0644))
	})
	a.EmailService = emailService
	a.EventsManager = events.New()
	a.Service = domain.New(a.Config.Service, a.Repo, a.Storage, a.TemplateService, a.EmailService)
	a.Server = httpserver.New(a.Config.Http, a.Service, a.AuthenticationService, a.EventsManager)
	assert.NoError(a.Server.Start())
	assert.NoError(initSessions(a.Sessions))

	httpexpectConfig := httpexpect.Config{
		BaseURL:  "http://" + a.Config.Http.Host + ":" + a.Config.Http.Port + "/",
		Reporter: httpexpect.NewRequireReporter(t),
	}
	if os.Getenv("DEBUG_HTTP") != "" {
		httpexpectConfig.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}
	e := httpexpect.WithConfig(httpexpectConfig)
	return e, func() {
		err := a.Shutdown()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}, sso
}

var UserSessions = make(map[int64]string)

func initSessions(sessionsStorage sessions.Sessions) error {
	for i := 1; i < 8; i++ {
		sessionToken, err := sessionsStorage.Add(int64(i))
		if err != nil {
			return err
		}
		UserSessions[int64(i)] = sessionToken
	}
	return nil
}

func InitDB(t *testing.T) (*require.Assertions, *app.Application) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.Mode = app.ModeTest
	config.Http.Host = "localhost"
	config.Http.Logger = false
	config.Templates.Dir = "../../templates/templates/"
	config.LogLevel = ""
	util.SetTime(util.Date("2019-04-01").Time)
	application := app.New(config)
	err := application.BootstrapDB()
	assert.NoError(err)

	return assert, application
}

func AuthInspecteur(request *httpexpect.Request) *httpexpect.Request {
	return AuthUser(request, 3)
}

func AuthExploitant(request *httpexpect.Request) *httpexpect.Request {
	return AuthUser(request, 1)
}

func AuthApprobateur(request *httpexpect.Request) *httpexpect.Request {
	return AuthUser(request, 6)
}

func AuthUser(request *httpexpect.Request, userId int64) *httpexpect.Request {
	return request.WithHeader(httpserver.AuthorizationHeader, UserSessions[userId])
}
