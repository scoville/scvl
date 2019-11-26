package engine

import (
	"fmt"
	"strings"

	"github.com/scoville/scvl/src/domain"
	"golang.org/x/crypto/bcrypt"
)

// FindUser finds and returns the user
func (e *Engine) FindUser(userID uint) (*domain.User, error) {
	return e.sqlClient.FindUser(&domain.User{ID: userID})
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

// RegistrationRequest is the request
type RegistrationRequest struct {
	Hash     string
	Password string
}

// UserRegister creates the user who is invited to the system.
func (e *Engine) UserRegister(req *RegistrationRequest) (*domain.User, error) {
	invitation, err := e.sqlClient.FindInvitation(req.Hash)
	if err != nil {
		return nil, err
	}
	if err := invitation.Valid(); err != nil {
		return nil, err
	}
	user, err := e.sqlClient.FindUser(&domain.User{
		Email:  invitation.ToUser.Email,
		Status: domain.UserStatusTemp,
	})
	if err != nil {
		return nil, err
	}
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}
	user.Status = domain.UserStatusValid
	usedInvitation, err := e.sqlClient.UpdateInvitation(invitation, &domain.UserInvitation{
		Status: domain.InvitationStatusUsed,
		ToUser: user,
	})
	return usedInvitation.ToUser, err
}

// LoginUserRequest is the Reqeust
type LoginUserRequest struct {
	Email    string
	Password string
}

// LoginUser is login request
func (e *Engine) LoginUser(req *LoginUserRequest) (*domain.User, error) {
	// todo: encrypt実装
	user, err := e.sqlClient.FindUser(&domain.User{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(req.Password)); err != nil {
		return nil, err
	}
	return user, nil
}
