package engine

import "github.com/scoville/scvl/src/domain"

// SQLClient is the interface for sql client
type SQLClient interface {
	FindUser(uint) (*domain.User, error)

	FindPageBySlug(string) (*domain.Page, error)
	CreatePage(*domain.Page) error
	UpdatePage(*domain.Page, *domain.Page) error

	CreatePageView(string, *domain.PageView) error

	FindOGPByID(int) (*domain.OGP, error)
	CreateOGP(ogp *domain.OGP) error
	UpdateOGP(uint, *domain.OGP) error
	DeleteOGP(uint) error

	FindFileBySlug(string) (*domain.File, error)
	CreateFile(*domain.File) error
	UpdateFile(*domain.File, domain.File) (err error)
}
