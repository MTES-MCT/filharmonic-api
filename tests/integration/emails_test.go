package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestSendEmailsRecapValidation(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	assert.NoError(application.Service.SendEmailsRecapValidation(int64(5)))
}
