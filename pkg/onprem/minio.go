package onprem

import (
	"context"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type storage struct {
	client     *minio.Client
	bucketName string
	urlPrefix  string
}

func NewStorage() (pkg.Storage, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT_URL")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	bucketPolicy := os.Getenv("MINIO_BUCKET_POLICY")
	useSSL := false // true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "minio.New")
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err == nil && !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			// Region: "us-east-1",
		})
		if err != nil {
			return nil, errors.Wrap(err, "minioClient.MakeBucket")
		}
		err = minioClient.SetBucketPolicy(ctx, bucketName, bucketPolicy)
		if err != nil {
			return nil, errors.Wrap(err, "minioClient.SetBucketPolicy")
		}
	}

	endpointURL := minioClient.EndpointURL()
	urlPrefix := endpointURL.Scheme + "://" + endpointURL.Host + "/" + bucketName + "/"

	return &storage{
		client:     minioClient,
		bucketName: bucketName,
		urlPrefix:  urlPrefix,
	}, nil
}

func (s *storage) Upload(ctx context.Context, input pkg.UploadInput) (string, error) {

	info, err := s.client.PutObject(ctx, s.bucketName, input.Name, input.File, input.Size,
		minio.PutObjectOptions{ContentType: input.ContentType},
	)
	if err != nil {
		return "", errors.Wrap(err, "s.client.PutObject")
	}

	return s.urlPrefix + info.Key, nil
}

func (s *storage) Delete(ctx context.Context, name string) error {
	name = strings.TrimPrefix(name, s.urlPrefix)
	err := s.client.RemoveObject(ctx, s.bucketName, name, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "s.client.RemoveObject")
	}
	return nil
}
