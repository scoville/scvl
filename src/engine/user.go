package engine

import (
	"fmt"
	"strings"

	"github.com/scoville/scvl/src/domain"
)

// FindUser finds and returns the user
func (e *Engine) FindUser(userID uint) (*domain.User, error) {
	return e.sqlClient.FindUser(userID)
}

// FindOrCreateUserByGoogleCode finds or creates the user
func (e *Engine) FindOrCreateUserByGoogleCode(code string) (*domain.User, error) {
	u, err := e.googleClient.FetchUserInfo(code)
	if err != nil {
		return nil, err
	}
	permitted := (e.allowedDomain == "")
	for _, allowedDomain := range strings.Split(e.allowedDomain, ",") {
		if strings.HasSuffix(u.Email, "@"+allowedDomain) {
			permitted = true
			break
		}
	}
	if !permitted {
		return nil, fmt.Errorf("only %s can allowed to use this service", e.allowedDomain)
	}
	return e.sqlClient.FindOrCreateUser(u)
}
