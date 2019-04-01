package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"

	// use postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type client struct {
	db *gorm.DB
}

// NewClient returns the engine.SQLClient
func NewClient(dbURL string) (engine.SQLClient, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(
		&domain.Email{},
		&domain.FileDownload{},
		&domain.File{},
		&domain.OGP{},
		&domain.PageView{},
		&domain.Page{},
		&domain.User{},
	)
	return &client{db}, nil
}

func (c *client) Close() error {
	return c.db.Close()
}
