package sql

const tblFiles = "files"

func (c *client) findFileBySlug(slug string) (file File, err error) {
	err = m.db.Table("files").
		First(&file, "slug = ?", slug).Error
	return
}

func (c *client) createFile(params *File) (err error) {
	err = m.db.Create(params).Error
	return
}

func (c *client) updateFile(file *File, params File) (err error) {
	return m.db.Table("files").
		Model(file).
		Update(&params).Error
}
