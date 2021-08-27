package engine

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/mssola/user_agent"
	"github.com/scoville/scvl/src/domain"
	"golang.org/x/crypto/bcrypt"
)

// FindFile fetches a file
func (e *Engine) FindFile(slug string, userID int) (file *domain.File, err error) {
	file, err = e.sqlClient.FindFileBySlug(slug)
	if err != nil {
		return
	}
	if userID != 0 && file.UserID != userID {
		err = errors.New("do not have permission")
		return
	}
	return
}

// UploadFileRequest is the request struct for the UploadFile function
type UploadFileRequest struct {
	DownloadLimit   int
	File            multipart.File
	FileName        string
	FileSize        int64
	Password        string
	UserID          int
	ValidDays       int
	SendEmail       bool
	DirectDownload  bool
	ReceiverAddress string
	ReceiverName    string
	SenderName      string
	BCCAddress      string
	Message         string
	SendPassword    bool
}

// UploadFile uploads a file to S3
func (e *Engine) UploadFile(req UploadFileRequest) (file *domain.File, err error) {
	slug := domain.GenerateSlug(44)
	path := fmt.Sprintf("%d/%s/%s", req.UserID, slug, req.FileName)
	err = e.awsClient.UploadToS3(req.File, path)
	if err != nil {
		return
	}

	var deadline *time.Time
	if req.ValidDays > 0 {
		t := time.Now().AddDate(0, 0, req.ValidDays)
		deadline = &t
	}

	encryptedPassword := ""
	if req.Password != "" {
		encryptedPassword = domain.Encrypt(req.Password)
	}

	file = &domain.File{
		UserID:            req.UserID,
		EncryptedPassword: encryptedPassword,
		Slug:              slug,
		Deadline:          deadline,
		Path:              path,
		DownloadLimit:     req.DownloadLimit,
		DirectDownload:    req.DirectDownload,
	}
	err = e.sqlClient.CreateFile(file)
	if err != nil {
		return
	}
	if !req.SendEmail {
		return
	}
	file.Email = &domain.FileEmail{
		FileID:          int(file.ID),
		SenderName:      req.SenderName,
		ReceiverName:    req.ReceiverName,
		ReceiverAddress: req.ReceiverAddress,
		BCCAddress:      req.BCCAddress,
		Message:         req.Message,
	}
	err = e.sqlClient.CreateFileEmail(file.Email)
	if err != nil {
		return
	}

	password := ""
	if req.SendPassword {
		password = req.Password
	}
	err = e.awsClient.SendFileEmail(file, password)
	return
}

// DownloadFileRequest is the request struct for DownloadFile function
type DownloadFileRequest struct {
	Slug      string
	Password  string
	RealIP    string
	Referer   string
	UserAgent string
}

// DownloadFile downloads a file from S3
func (e *Engine) DownloadFile(req *DownloadFileRequest) (fileName string, data []byte, err error) {
	file, err := e.sqlClient.FindFileBySlug(req.Slug)
	if err != nil {
		return
	}
	if err = file.Downloadable(); err != nil {
		return
	}

	if file.EncryptedPassword != "" {
		err = bcrypt.CompareHashAndPassword(
			[]byte(file.EncryptedPassword),
			[]byte(req.Password))
		if err != nil {
			return
		}
	}
	fileName = filepath.Base(file.Path)
	data, err = e.awsClient.DownloadFromS3(file.Path)
	if err != nil {
		return
	}
	ua := user_agent.New(req.UserAgent)
	name, _ := ua.Browser()
	e.sqlClient.CreateFileDownload(req.Slug, &domain.FileDownload{
		RealIP:      req.RealIP,
		Referer:     req.Referer,
		Mobile:      ua.Mobile(),
		Platform:    ua.Platform(),
		OS:          ua.OS(),
		BrowserName: name,
	})
	return
}

// UpdateFileRequest is the request struct for UpdateFile function
type UpdateFileRequest struct {
	Slug          string
	Status        string
	DownloadLimit int
	File          multipart.File
	FileName      string
	FileSize      int64
	Password      string
	UserID        int
}

// UpdateFile updates the file
func (e *Engine) UpdateFile(req *UpdateFileRequest) (file *domain.File, err error) {
	file, err = e.sqlClient.FindFileBySlug(req.Slug)
	if err != nil {
		return
	}

	if file.UserID != req.UserID {
		err = errors.New("You don't have permission to edit it")
		return
	}

	path := file.Path
	if req.FileName != "" {
		path = fmt.Sprintf("%d/%s/%s", req.UserID, req.Slug, req.FileName)
		err = e.awsClient.UploadToS3(req.File, path)
		if err != nil {
			return
		}
	}

	encryptedPassword := ""
	if req.Password != "" {
		encryptedPassword = domain.Encrypt(req.Password)
	}
	status := req.Status
	if status == "" {
		status = file.Status
	}

	err = e.sqlClient.UpdateFile(file, &domain.File{
		EncryptedPassword: encryptedPassword,
		Path:              path,
		DownloadLimit:     req.DownloadLimit,
		Status:            status,
	})
	return
}
