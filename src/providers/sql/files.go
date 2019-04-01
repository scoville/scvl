package sql

import "github.com/scoville/scvl/src/domain"

const tblFiles = "files"

func (c *client) FindFileBySlug(slug string) (file *domain.File, err error) {
	file = &domain.File{}
	err = c.db.Table(tblFiles).
		First(&file, "slug = ?", slug).Error
	if err != nil {
		return
	}
	err = c.db.Table(tblFileDownloads).
		Where("file_id = ?", file.ID).
		Count(&(file.DownloadCount)).Error
	return
}

func (c *client) CreateFile(params *domain.File) (err error) {
	err = c.db.Create(params).Error
	return
}

func (c *client) UpdateFile(file, params *domain.File) (err error) {
	return c.db.Table(tblFiles).
		Model(file).
		Update(&params).Error
}
