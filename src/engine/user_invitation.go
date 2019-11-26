package engine

import (
	"github.com/scoville/scvl/src/domain"
)

// FindInvitationRequest is the Request
type FindInvitationRequest struct {
	Hash string
}

// FindInvitation find the user invitation by hash
func (e *Engine) FindInvitation(req *FindInvitationRequest) (*domain.UserInvitation, error) {
	return e.sqlClient.FindInvitation(req.Hash)
}

// InviteRequest is the request
type InviteRequest struct {
	FromUserID uint
	Email      string
}

// InviteUser deals new user which is invited by existing user
func (e *Engine) InviteUser(req *InviteRequest) (*domain.UserInvitation, error) {
	if _, err := e.sqlClient.FindUser(&domain.User{ID: req.FromUserID}); err != nil {
		return nil, err
	}
	paramas := &domain.UserInvitation{
		FromUserID: req.FromUserID,
		ToUser: &domain.User{
			Status: domain.UserStatusTemp,
			Email:  req.Email,
		},
	}
	invitation, err := e.sqlClient.CreateInvitation(paramas)
	return invitation, err
}
