package storage

import (
	"errors"
	"io"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gofrs/uuid"
	minio "github.com/minio/minio-go"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Endpoint        string `default:"localhost:9000"`
	UseSSL          bool   `default:"false"`
	AccessKey       string `default:"filharmonic"`
	SecretAccessKey string `default:"filharmonic"`
	BucketName      string `default:"filharmonic"`
}

type FileStorage struct {
	config Config

	client *minio.Client
}

func New(config Config) (*FileStorage, error) {
	client, err := minio.New(config.Endpoint, config.AccessKey, config.SecretAccessKey, config.UseSSL)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("connecting to minio endpoint on %s", config.Endpoint)
	exists, err := client.BucketExists(config.BucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		log.Warn().Msg("missing minio bucket, creating...")
		err = client.MakeBucket(config.BucketName, "")
		if err != nil {
			return nil, err
		}
	}

	return &FileStorage{
		config: config,
		client: client,
	}, nil
}

func (storage *FileStorage) Put(pieceJointe models.PieceJointeFile) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	idStr := id.String()
	available, err := storage.checkKeyAvailable(idStr)
	if err != nil {
		return "", err
	}
	if !available {
		return "", errors.New("existing object at key " + idStr)
	}
	_, err = storage.client.PutObject(storage.config.BucketName, idStr, pieceJointe.Content, pieceJointe.Taille, minio.PutObjectOptions{})
	return idStr, err
}

func (storage *FileStorage) Get(id string) (io.Reader, error) {
	available, err := storage.checkKeyAvailable(id)
	if err != nil {
		return nil, err
	}
	if available {
		return nil, errors.New("missing object at key " + id)
	}
	obj, err := storage.client.GetObject(storage.config.BucketName, id, minio.GetObjectOptions{})
	return obj, err
}

func (storage *FileStorage) checkKeyAvailable(id string) (bool, error) {
	_, err := storage.client.StatObject(storage.config.BucketName, id, minio.StatObjectOptions{})
	if err == nil {
		return false, nil
	}
	switch v := err.(type) {
	case minio.ErrorResponse:
		if v.Code == "NoSuchKey" {
			return true, nil
		}
		return false, err
	default:
		return false, err
	}
}
