package sql

import "github.com/scoville/scvl/src/domain"

func (c *client) CreateEmail(email *domain.Email) (err error) {
	return c.db.Create(email).Error
}
