package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/scoville/scvl/src/domain"
)

const tblUsers = "users"

func (c *client) FindUser(id uint) (user *domain.User, err error) {
	user = &domain.User{}

	err = c.db.Table(tblUsers).
		Preload("Pages", func(db *gorm.DB) *gorm.DB {
			return db.Order("pages.created_at DESC")
		}).
		Preload("Pages.OGP").
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			return db.Order("files.created_at DESC")
		}).
		First(user, id).Error
	if err != nil {
		return
	}
	for _, p := range user.Pages {
		c.db.Table(tblPageviews).Where("page_id = ?", p.ID).Count(&(p.ViewCount))
	}
	for _, f := range user.Files {
		c.db.Table(tblFileDownloads).Where("file_id = ?", f.ID).Count(&(f.DownloadCount))
	}

	return
}

func (c *client) findOrCreateUser(params User) (user User, err error) {
	err = c.db.
		Where(User{Email: params.Email}).
		Assign(User{Name: params.Name}).
		FirstOrCreate(&user).Error
	return
}
