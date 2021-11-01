package storage

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/grigagod/compresso/pkg/db/aws"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAWSStorage_BasicOperations(t *testing.T) {
	client, err := aws.NewClientWithEnvCredentials()
	assert.NoError(t, err)

	cfg := &Config{
		Bucket:          "compresso-test",
		PresignDuration: 1 * time.Second,
	}
	storage := NewAWSStorage(cfg, client)
	file := []byte(string("hello S3!"))

	t.Run("PutObject", func(t *testing.T) {
		err := storage.PutObject(context.Background(), bytes.NewReader(file), string("test_basic.txt"))
		assert.NoError(t, err)
	})

	t.Run("GetObject", func(t *testing.T) {
		resp, err := storage.GetObject(context.Background(), "test_basic.txt")
		defer resp.Close()
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(resp)
		assert.NoError(t, err)
		assert.Equal(t, body, file)
	})

	t.Run("GetDownloadURL", func(t *testing.T) {
		url, err := storage.GetDownloadURL("test_basic.txt")
		assert.NoError(t, err)

		resp, err := http.Get(url)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, body, file)
	})

	t.Run("PutObjectContextTimeout", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), time.Microsecond)
		err := storage.PutObject(ctx, bytes.NewReader(file), string("test_basic.txt"))
		assert.Error(t, err)

		awsErr, ok := errors.Cause(err).(awserr.Error)
		assert.True(t, ok)
		assert.Equal(t, awsErr.Code(), request.CanceledErrorCode)
	})

	t.Run("GetObjectContextTimeout", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), time.Microsecond)
		_, err := storage.GetObject(ctx, "test_basic.txt")
		assert.Error(t, err)

		awsErr, ok := errors.Cause(err).(awserr.Error)
		assert.True(t, ok)
		assert.Equal(t, awsErr.Code(), request.CanceledErrorCode)

	})

	t.Run("GetDownloadURLPresignExpired", func(t *testing.T) {
		url, err := storage.GetDownloadURL("test_basic.txt")
		assert.NoError(t, err)

		time.Sleep(time.Second)

		resp, err := http.Get(url)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NotEqual(t, body, file)
	})
}
