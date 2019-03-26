package engine

import (
	"mime/multipart"
	"time"

	"github.com/scoville/scvl/src/domain"
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

// UploadFile uploads file to S3
func (e *Engine) UploadFile(req UploadFileRequest) (file *domain.File, err error) {
	err = s3Client.UploadToS3(file, "test.jpg")
	if err != nil {
		return
	}

	var deadline *time.Time
	if req.ValidDays > 0 {
		t := time.Now().AddDate(0, 0, req.ValidDays)
		deadline = &t
	}

	err = manager.createFile(&domain.File{
		UserID:            int(user.ID),
		EncryptedPassword: domain.Encrypt(password),
		Slug:              domain.GenerateSlug(44),
		Deadline:          deadline,
	})
	if err != nil {
		return
	}

}
