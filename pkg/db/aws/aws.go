// Package aws provides connection to AWS S3 service.
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewClientWithStaticCredentials create new S3 client with static credentials.
func NewClientWithStaticCredentials(region, acccessKeyID, secretAccessKey string) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(acccessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

// NewClientWithSharedCredentials create new S3 client with shared credentials.
func NewClientWithSharedCredentials(filename, profile string) (*s3.S3, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigFiles: []string{filename},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}
