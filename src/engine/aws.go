package engine

import "io"

// AWSClient is the interface for aws client
type AWSClient interface {
	UploadToS3(io.ReadSeeker, string) error
	DownloadFromS3(string) ([]byte, error)
	SendMail(string, string, string) error
}
