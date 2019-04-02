package awsclient

import (
	"bytes"
	"errors"
	"io"
	"net/mail"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

type awsClient struct {
	svc              *session.Session
	s3Bucket         string
	s3Region         string
	sesRegion        string
	mailBccAddresses []*string
	mailFrom         mail.Address
	baseURL          string
}

// Config is params for NewClient
type Config struct {
	AccessKey      string
	AccessSecret   string
	S3Bucket       string
	S3Region       string
	SESRegion      string
	MailFrom       string
	MailBCCAddress string
	BaseURL        string
}

// NewClient creates and returns awsClient
func NewClient(config Config) (engine.AWSClient, error) {
	svc, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			config.AccessKey,
			config.AccessSecret,
			"",
		),
	})
	if err != nil {
		return nil, err
	}

	return &awsClient{
		svc:       svc,
		s3Bucket:  config.S3Bucket,
		s3Region:  config.S3Region,
		sesRegion: config.SESRegion,
		mailFrom: mail.Address{
			Name:    "scvl",
			Address: config.MailFrom,
		},
		mailBccAddresses: []*string{aws.String(config.MailBCCAddress)},
		baseURL:          config.BaseURL,
	}, nil
}

func (c *awsClient) UploadToS3(file io.ReadSeeker, path string) error {
	_, err := s3.New(c.svc, aws.NewConfig().WithRegion(c.s3Region)).
		PutObject(&s3.PutObjectInput{
			Bucket: aws.String(c.s3Bucket),
			Key:    aws.String(path),
			Body:   file,
		})
	return err
}

func (c *awsClient) DownloadFromS3(path string) (data []byte, err error) {
	resp, err := s3.New(c.svc, aws.NewConfig().WithRegion(c.s3Region)).
		GetObject(&s3.GetObjectInput{
			Bucket: aws.String(c.s3Bucket),
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

func (c *awsClient) SendMail(file *domain.File, password string) error {
	if file.Email == nil {
		return errors.New("Email is empty")
	}
	toAddresses := []*string{}
	for _, email := range strings.Split(file.Email.ReceiverAddress, ",") {
		if email != "" {
			toAddresses = append(toAddresses, aws.String(strings.TrimSpace(email)))
		}
	}
	bccAddresses := c.mailBccAddresses
	for _, email := range strings.Split(file.Email.BCCAddress, ",") {
		if email != "" {
			bccAddresses = append(bccAddresses, aws.String(strings.TrimSpace(email)))
		}
	}
	body, err := file.Email.HTML(c.baseURL, file)
	if err != nil {
		return err
	}

	svc := ses.New(c.svc, aws.NewConfig().WithRegion(c.sesRegion))
	_, err = svc.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  toAddresses,
			BccAddresses: bccAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("ファイルが届いています。"),
			},
		},
		Source: aws.String(c.mailFrom.String()),
	})
	if err != nil {
		return err
	}
	if password == "" {
		return nil
	}
	passBody, err := file.Email.PasswordHTML(file, password)
	if err != nil {
		return err
	}
	_, err = svc.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  toAddresses,
			BccAddresses: bccAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(passBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("パスワードをお知らせします"),
			},
		},
		Source: aws.String(c.mailFrom.String()),
	})
	return err
}
