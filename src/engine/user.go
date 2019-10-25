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
	var ok bool
	for _, allowedDomain := range strings.Split(e.allowedDomains, ",") {
		if strings.HasSuffix(u.Email, "@"+allowedDomain) {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("this email address is not allowed to use this service")
	}

	return e.sqlClient.FindOrCreateUser(u)
}
