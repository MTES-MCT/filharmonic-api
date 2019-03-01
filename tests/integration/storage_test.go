package integration

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	assert := require.New(t)

	config := storage.Config{
		Endpoint:        "localhost:9000",
		AccessKey:       "filharmonic",
		SecretAccessKey: "filharmonic",
		BucketName:      "filharmonic",
	}
	store, err := storage.New(config)
	assert.NoError(err)

	content := "MonContenu"
	id, err := store.Put(models.File{
		Content: strings.NewReader(content),
		Type:    "application/pdf",
		Taille:  int64(len(content)),
		Nom:     "test.pdf",
	})
	assert.NoError(err)
	assert.NotEmpty(id)

	reader, err := store.Get(id)
	assert.NoError(err)
	data, err := ioutil.ReadAll(reader)
	assert.NoError(err)
	assert.Equal(content, string(data))
}
