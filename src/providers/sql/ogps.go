package sql

import "github.com/scoville/scvl/src/domain"

func (c *client) FindOGPByID(id int) (ogp *domain.OGP, err error) {
	ogp = &domain.OGP{}
	err = c.db.Table("ogps").
		Where("id = ?", id).
		First(&ogp).Error
	return
}

func (c *client) CreateOGP(ogp *domain.OGP) (err error) {
	return c.db.Create(ogp).Error
}

func (c *client) UpdateOGP(id uint, ogp *domain.OGP) (err error) {
	return c.db.Table("ogps").
		Where("id = ?", id).
		Update(ogp).Error
}

func (c *client) DeleteOGP(id uint) (err error) {
	return c.db.Table("ogps").
		Delete(domain.OGP{}, "id = ?", id).Error
}
