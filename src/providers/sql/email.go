package sql

import (
	"time"

	"github.com/scoville/scvl/src/domain"
)

const tblEmails = "emails"

func (c *client) CreateEmail(email *domain.Email) (err error) {
	err = c.db.Create(email).Error
	return
}

func (c *client) ReadEmail(emailID string) (err error) {
	now := time.Now()
	return c.db.Table(tblEmails).
		Where("id = ?", emailID).
		UpdateColumn("opened_at", now).Error
}
