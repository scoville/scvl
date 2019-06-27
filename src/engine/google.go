package engine

import "github.com/scoville/scvl/src/domain"

// GoogleClient is the interface for google client
type GoogleClient interface {
	FetchUserInfo(string) (domain.User, error)
	GetDriveFileTitle(user *domain.User, id string) (title string, err error)
	AuthCodeURL(string) string
}
