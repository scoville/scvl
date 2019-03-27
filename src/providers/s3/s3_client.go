package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Client struct {
	S3Bucket string
	svc      *session.Session
}

func newS3Client(s3Bucket, s3Region, id, secret string) (engine.S3Client, error) {
	svc, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
		Region:      aws.String(s3Region),
	})
	if err != nil {
		return nil, err
	}

	return &s3Client{
		S3Bucket: s3Bucket,
		svc:      svc,
	}, nil
}

func (c *s3Client) Upload(file io.ReadSeeker, path string) error {
	_, err := s3.New(c.svc).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(c.S3Bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	return err
}

func (c *s3Client) Download(path string) (data []byte, err error) {
	resp, err := s3.New(c.svc).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.S3Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	data = buf.Bytes()
	return
}
