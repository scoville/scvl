package engine

import (
	"errors"
	"strings"

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
	for _, row := range sheet.Rows {
		for i, col := range row {
			variables[keys[i]] = col.Value
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
	return
}
