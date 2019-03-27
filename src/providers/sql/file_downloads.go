package sql

import "github.com/scoville/scvl/src/domain"

const tblFileDownloads = "file_downloads"

func (c *client) createFileDownload(slug string, fd domain.FileDownload) (err error) {
	var file domain.File
	err = c.db.Where(&domain.File{Slug: slug}).First(&file).Error
	if err != nil {
		return
	}
	fd.FileID = int(file.ID)
	return c.db.Create(&fd).Error
}
