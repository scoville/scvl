package engine

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/scoville/scvl/src/domain"
	"golang.org/x/crypto/bcrypt"
)

// UploadFileRequest is the request struct for the UploadFile function
type UploadFileRequest struct {
	DownloadLimit int
	Email         string
	File          multipart.File
	FileName      string
	FileSize      int64
	Password      string
	UserID        int
	ValidDays     int
}

// UploadFile uploads a file to S3
func (e *Engine) UploadFile(req UploadFileRequest) (file *domain.File, err error) {
	slug := domain.GenerateSlug(44)
	path := fmt.Sprintf("%d/%s/%s", req.UserID, slug, req.FileName)
	err = e.s3Client.Upload(req.File, path)
	if err != nil {
		return
	}

	var deadline *time.Time
	if req.ValidDays > 0 {
		t := time.Now().AddDate(0, 0, req.ValidDays)
		deadline = &t
	}

	file = &domain.File{
		UserID:            req.UserID,
		EncryptedPassword: domain.Encrypt(req.Password),
		Slug:              slug,
		Deadline:          deadline,
		Path:              path,
		DownloadLimit:     req.DownloadLimit,
	}
	err = e.sqlClient.CreateFile(file)
	return
}

// DownloadFile downloads a file from S3
func (e *Engine) DownloadFile(slug, password string) (data []byte, err error) {
	file, err := e.sqlClient.FindFileBySlug(slug)
	if err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(file.EncryptedPassword),
		[]byte(password)); err != nil {
		return
	}

	return e.s3Client.Download(file.Path)
}
