package engine

import "github.com/scoville/scvl/src/domain"

// FindUser finds and returns the user
func (e *Engine) FindUser(userID uint) (u *domain.User, err error) {
	return e.sqlClient.FindUser(userID)
}
