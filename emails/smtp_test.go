package emails

// func TestSend(t *testing.T) {
// 	assert := require.New(t)

// 	emailConfig := Config{
// 		FromName: "Fil'Harmonic",
// 		SmtpHost: "mail.filharmonic.beta.gouv.fr",
// 		SmtpPort: 25,
// 		SmtpUser: "noreply@filharmonic.beta.gouv.fr",
// 		SmtpPass: "", // besoin mdp
// 	}
// 	emailService, err := New(emailConfig)
// 	assert.NoError(err)

// 	err = emailService.Send(Email{
// 		To:       "replace-me@filharmonic.beta.gouv.fr",
// 		Subject:  "Test Email",
// 		TextPart: "message text",
// 		HTMLPart: "<b>message html</b>",
// 	})
// 	assert.NoError(err)
// }
