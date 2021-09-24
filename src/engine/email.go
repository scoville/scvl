package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scoville/scvl/src/domain"
)

// CreateEmailRequest is the request for creating an email
type CreateEmailRequest struct {
	IsPreview      bool
	SpreadsheetURL string
	SheetName      string
	Sender         string
	Template       string
	Title          string
	User           *domain.User
}

// CreateEmail creates an email
func (e *Engine) CreateEmail(req *CreateEmailRequest) (emailTemplate *domain.EmailTemplate, err error) {

	emailTemplate, err = e.createEmailTemplate(req)
	return
}

// SendEmail send a email
func (e *Engine) SendEmail(req *CreateEmailRequest) (err error) {
	// base64でエンコード済みの透過gifを本文に追加
	emailTemplate, err := e.createEmailTemplate(req)
	if err != nil {
		return
	}
	for _, email := range emailTemplate.BatchEmail.Emails {
		openConfirmationCode := fmt.Sprintf(`<img src="https://%s/%d/read">`, os.Getenv("EMAIL_DOMAIN"), email.ID)
		email.Body = email.Body + openConfirmationCode
		err = e.awsClient.SendEmail(email, emailTemplate.BatchEmail.Sender)
		if err != nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
	return
}

// ReadEmail update emails opened_at
func (e *Engine) ReadEmail(emailID string) (err error) {
	err = e.sqlClient.ReadEmail(emailID)
	return
}

func (e *Engine) createEmailTemplate(req *CreateEmailRequest) (emailTemplate *domain.EmailTemplate, err error) {
	splitted := strings.Split(req.SpreadsheetURL, "/")
	if len(splitted) < 6 {
		err = errors.New("invalid spreadsheet url")
		return
	}

	spreadsheet, err := e.googleClient.FetchSpreadsheet(req.User, splitted[5])
	if err != nil {
		return
	}
	sheet, err := spreadsheet.SheetByTitle(req.SheetName)
	if err != nil {
		return
	}
	if len(sheet.Rows) < 2 {
		err = errors.New("送信するメールがシート内に見つかりませんでした")
		return
	}

	emailTemplate = &domain.EmailTemplate{
		UserID: int(req.User.ID),
		Body:   req.Template,
		Title:  req.Title,
		BatchEmail: &domain.BatchEmail{
			Sender:         req.Sender,
			SpreadsheetURL: req.SpreadsheetURL,
			SheetName:      req.SheetName,
		},
	}

	keys := make([]string, 0, len(sheet.Rows[0]))
	for _, col := range sheet.Rows[0] {
		keys = append(keys, col.Value)
	}
	variables := make(map[string]string, len(keys))
	for i, row := range sheet.Rows {
		if i == 0 {
			continue
		}
		for t, col := range row {
			variables[keys[t]] = col.Value
		}
		to, ok := variables["email"]
		if !ok {
			err = errors.New("シートから送信先のemailが見つかりませんでした")
			return
		}
		var email *domain.Email
		email, err = domain.NewEmail(req.Template, to, req.Title, variables)
		if err != nil {
			return
		}
		emailTemplate.BatchEmail.Emails = append(emailTemplate.BatchEmail.Emails, email)
	}
	emailTemplate.BatchEmail.SentCount = len(emailTemplate.BatchEmail.Emails)
	err = e.sqlClient.CreateEmailTemplate(emailTemplate)
	return emailTemplate, err
}
