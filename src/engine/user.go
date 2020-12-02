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
	permitted := (e.allowedDomains == "")
	for _, allowedDomain := range strings.Split(e.allowedDomains, ",") {
		if strings.HasSuffix(u.Email, "@"+allowedDomain) {
			permitted = true
			break
		}
	}
	if !permitted {
		return nil, fmt.Errorf("this email address is not allowed to use this service")
	}

	return e.sqlClient.FindOrCreateUser(u)
}

// UpdateUserAPIKey updates the user
func (e *Engine) UpdateUserAPIKey(userID uint) error {
	u, err := e.sqlClient.FindUser(userID)
	if err != nil {
		return err
	}
	key := domain.GenerateSlug(40)
	return e.sqlClient.UpdateUser(u, &domain.User{APIKey: key})
}
