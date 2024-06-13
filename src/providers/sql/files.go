package sql

import (
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

const tblFiles = "files"

func (c *client) FindFiles(params *engine.FindFilesRequest) (files []*domain.File, count int, err error) {
	db := c.db.Table(tblFiles).
		Where("user_id = ?", params.UserID).
		Where("status = ?", domain.FileStatusActive)

	if params.Query != "" {
		db = db.Where("path LIKE ?", "%" + params.Query + "%")
	}

	db = db.Count(&count)

	if params.Limit != 0 {
		db = db.Limit(params.Limit)
	}
	if params.Offset != 0 {
		db = db.Offset(params.Offset)
	}

	err = db.Order("created_at DESC").
		Find(&files).Error
	if err != nil {
		return
	}
	for _, file := range files {
		c.db.Table(tblFileDownloads).Where("file_id = ?", file.ID).Count(&(file.DownloadCount))
	}

	return
}

func (c *client) FindFileBySlug(slug string) (file *domain.File, err error) {
	file = &domain.File{}
	err = c.db.Table(tblFiles).
		Where("status = ?", domain.FileStatusActive).
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
