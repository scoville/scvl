package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Manager はデータベースとやり取りをします
type Manager struct {
	db *gorm.DB
}

var manager *Manager

func setupManager() {
	db, err := gorm.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err.Error())
	}
	db.AutoMigrate(&User{}, &Page{}, &PageView{})
	manager = &Manager{db}
}

func (m *Manager) findUser(id uint) (user User, err error) {
	err = m.db.First(&user, id).Error
	return
}

func (m *Manager) findOrCreateUser(u User) (user User, err error) {
	err = m.db.
		Where(User{Email: u.Email}).
		Assign(User{Name: u.Name}).
		FirstOrCreate(&user).Error
	return
}

func (m *Manager) setPagesToUser(u *User) (err error) {
	var pages []*Page
	err = m.db.
		Where(Page{UserID: int(u.ID)}).
		Order("created_at DESC").
		Find(&pages).Error
	if err != nil {
		return
	}
	for _, p := range pages {
		m.db.Model(&PageView{}).Where("page_id = ?", p.ID).Count(&(p.ViewCount))
	}
	u.Pages = pages
	return
}

func (m *Manager) createPage(userID uint, slug, url string) (err error) {
	return m.db.Create(&Page{
		UserID: int(userID),
		Slug:   slug,
		URL:    url,
	}).Error
}

func (m *Manager) createPageView(slug string, pv PageView) (err error) {
	var page Page
	err = m.db.Where(&Page{Slug: slug}).First(&page).Error
	if err != nil {
		return
	}
	pv.PageID = int(page.ID)
	return m.db.Create(&pv).Error
}
