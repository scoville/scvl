package sql

import "github.com/scoville/scvl/src/domain"

const tblImages = "images"

func (c *client) CreateImage(params *domain.Image) (err error) {
	err = c.db.Create(params).Error
	return
}
