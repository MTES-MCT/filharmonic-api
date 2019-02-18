package emails

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	assert := require.New(t)

	// start a fake MailJet API
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"success": true}`))
		})
		http.ListenAndServe(":5100", nil)
	}()

	emailConfig := Config{
		FromEmail:     "contact@filharmonic.beta.gouv.fr",
		FromName:      "Fil'Harmonic",
		APIPublicKey:  "publickey",
		APIPrivateKey: "privatekey",
	}
	emailService := New(emailConfig)
	emailService.SetBaseURL("http://localhost:5100/")

	err := emailService.Send(Email{
		Subject:        "Test Email",
		RecipientEmail: "sample-user@filharmonic.beta.gouv.fr",
		RecipientName:  "Sample user",
		TextPart:       "message",
		HTMLPart:       "<b>message</b>",
	})
	assert.NoError(err)
}
