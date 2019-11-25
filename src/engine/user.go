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
	if e.allowedDomain != "" && !strings.HasSuffix(u.Email, "@"+e.allowedDomain) {
		return nil, fmt.Errorf("only %s can allowed to use this service", e.allowedDomain)
	}
	return e.sqlClient.FindOrCreateUser(u)
}

// InviteRequest is the request
type InviteRequest struct {
	FromUserID uint
	Email      string
}

// InviteUser deals new user which is invited by existing user
func (e *Engine) InviteUser(req *InviteRequest) (*domain.UserInvitation, error) {
	if _, err := e.sqlClient.FindUser(req.FromUserID); err != nil {
		return nil, err
	}
	paramas := &domain.UserInvitation{
		FromUserID: req.FromUserID,
		ToUser: &domain.User{
			Email: req.Email,
		},
	}
	invitation, err := e.sqlClient.CreateInvitation(paramas)
	return invitation, err
}

// RegistrationRequest is the request
type RegistrationRequest struct {
	Hash     string
	Email    string
	Password string
}

// CreateUser creates the user who is invited to the system.
func (e *Engine) UserRegister(req *RegistrationRequest) (*domain.User, error) {
	if _, err := e.sqlClient.FindInvitation(req.Hash); err != nil {
		return nil, err
	}
	// Todo: invitationを使用済みにする
	paramas := &domain.User{
		// Password: req.Password, // passwordを追加するかどうか問題。
	}
	user, err := e.sqlClient.UserRegister(paramas)
	return user, err
}
