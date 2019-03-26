package sql

func (c *client) findOGPByID(id int) (ogp *OGP, err error) {
	ogp = &OGP{}
	err = m.db.Table("ogps").
		Where("id = ?", id).
		First(&ogp).Error
	return
}

func (c *client) createOGP(ogp *OGP) (err error) {
	return m.db.Create(ogp).Error
}

func (c *client) updateOGP(id uint, ogp OGP) (err error) {
	return m.db.Table("ogps").
		Where("id = ?", id).
		Update(&ogp).Error
}

func (c *client) deleteOGP(id uint) (err error) {
	return m.db.Table("ogps").
		Delete(OGP{}, "id = ?", id).Error
}
