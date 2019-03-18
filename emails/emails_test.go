package emails

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailsToBytes(t *testing.T) {
	assert := require.New(t)
	content, err := (&Email{
		From:    "Fil'harmonic <noreply@filharmonic.beta.gouv.fr>",
		To:      "sample-user@filharmonic.beta.gouv.fr",
		Subject: "Fil'Harmonic : Vous avez des nouveaux messages",
		TextPart: `Bonjour,

		Vous avez des nouveaux messages sur https://filharmonic.beta.gouv.fr

		A bientôt sur Fil'Harmonic
		`,
		HTMLPart: `<!DOCTYPE html><html><meta charset="UTF-8"><body><h1>Bonjour,</h1>

		<p>Vous avez des nouveaux messages sur <a href="https://filharmonic.beta.gouv.fr">filharmonic.beta.gouv.fr</a></p>

		<p>A bientôt sur Fil'Harmonic</p>
		</body></html>`,
	}).ToBytes()
	assert.NoError(err)
	assert.Contains(string(content), "Subject: Fil'Harmonic : Vous avez des nouveaux messages")
	assert.Contains(string(content), "A bient=C3=B4t sur Fil'Harmonic")
}
