package sql

import (
	"time"

	"github.com/scoville/scvl/src/domain"
)

const tblEmails = "emails"

func (c *client) CreateEmailTemplate(emailTemplate *domain.EmailTemplate) (err error) {
	err = c.db.Create(emailTemplate).Error
	return
}

func (c *client) ReadEmail(emailID string) (err error) {
	email := &domain.Email{}
	err = c.db.Table(tblEmails).
		Find(&email, "id = ?", emailID).Error
	if err != nil || email.OpenedAt != nil {
		return
	}
	now := time.Now()
	email.OpenedAt = &now
	err = c.db.Save(&email).Error
	if err != nil {
		return
	}

	batchEmail := domain.BatchEmail{}
	err = c.db.Table("batch_emails").
		Find(&batchEmail, "id = ?", email.BatchEmailID).Error
	if err != nil {
		return
	}
	batchEmail.OpenCount += 1
	err = c.db.Save(&batchEmail).Error
	return
}
