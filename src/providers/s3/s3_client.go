package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Client interacts with s3 to download/upload files
type S3Client struct {
	S3Bucket string
	svc      *session.Session
}

func newS3Client(s3Bucket, s3Region, id, secret string) (*S3Client, error) {
	svc, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
		Region:      aws.String(s3Region),
	})
	if err != nil {
		return nil, err
	}

	return &S3Client{
		S3Bucket: s3Bucket,
		svc:      svc,
	}, nil
}

// UploadToS3 uploads a image to S3
func (c *S3Client) UploadToS3(file io.ReadSeeker, path string) error {
	_, err := s3.New(c.svc).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(c.S3Bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	return err
}

// FetchFromS3 fetches a image from S3
func (c *S3Client) FetchFromS3(name string) (data []byte, err error) {
	resp, err := s3.New(c.svc).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.S3Bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	data = buf.Bytes()
	return
}
