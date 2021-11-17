package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// AWSStorage imlements Storage interface.
type AWSStorage struct {
	cfg      *Config
	client   *s3.S3
	uploader *s3manager.Uploader
}

// NewAWSStorage create new AWSStorage with given config and S3 client.
func NewAWSStorage(cfg *Config, client *s3.S3) *AWSStorage {
	return &AWSStorage{
		cfg:      cfg,
		client:   client,
		uploader: s3manager.NewUploaderWithClient(client),
	}
}

// PutObject upoad given file to the bucket.
func (s *AWSStorage) PutObject(ctx context.Context, file io.Reader, fileName, fileType string) error {
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return errors.Wrap(err, "AWSStorage.PutObject")
	}

	return nil
}

// GetObject return io.ReadCloser for given file from the bucket.
func (s *AWSStorage) GetObject(ctx context.Context, fileName string) (io.ReadCloser, error) {
	resp, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "AWSStorage.GetObject.GetObjectWithContext")
	}

	return resp.Body, nil
}

// GetDownloadURL create Presigned URL for downloading given file.
func (s *AWSStorage) GetDownloadURL(fileName string) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(fileName),
	})
	url, err := req.Presign(s.cfg.PresignDuration)
	if err != nil {
		return "", errors.Wrap(err, "AWSStorage.GetDownloadURL.Presign")
	}

	return url, nil
}
