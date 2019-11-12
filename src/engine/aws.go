package engine

import (
	"io"

	"github.com/scoville/scvl/src/domain"
)

// AWSClient is the interface for aws client
type AWSClient interface {
	UploadToS3(io.ReadSeeker, string) error
	DownloadFromS3(string) ([]byte, error)
	SendFileMail(*domain.File, string) error
	SendGroupMails(toAddresses, bccAdresses []*string, body string) error
}
