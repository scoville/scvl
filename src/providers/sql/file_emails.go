package sql

import "github.com/scoville/scvl/src/domain"

func (c *client) CreateFileEmail(email *domain.FileEmail) (err error) {
	return c.db.Create(email).Error
}
