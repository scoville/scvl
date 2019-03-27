package engine

import "io"

// S3Client is the interface for s3 client
type S3Client interface {
	Upload(io.ReadSeeker, string) error
	Download(string) ([]byte, error)
}
