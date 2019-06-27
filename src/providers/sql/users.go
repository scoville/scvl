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
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("images.created_at DESC")
		}).
		First(user, id).Error
	if err != nil {
		return
	}
	for _, p := range user.Pages {
		c.db.Table(tblPageViews).Where("page_id = ?", p.ID).Count(&(p.ViewCount))
	}
	for _, f := range user.Files {
		c.db.Table(tblFileDownloads).Where("file_id = ?", f.ID).Count(&(f.DownloadCount))
	}

	return
}

func (c *client) FindOrCreateUser(params domain.User) (user *domain.User, err error) {
	user = &domain.User{}
	err = c.db.
		Where(domain.User{Email: params.Email}).
		Assign(domain.User{Name: params.Name, GoogleToken: params.GoogleToken}).
		FirstOrCreate(user).Error
	return
}
