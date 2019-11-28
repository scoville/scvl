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
	invitation, err := e.sqlClient.FindInvitation(domain.UserInvitation{Hash: req.Hash})
	if err != nil {
		return nil, err
	}
	err = invitation.Valid()
	return invitation, err
}

// InviteRequest is the request
type InviteRequest struct {
	FromUserID uint
	Email      string
}

// InviteUser deals new user which is invited by existing user
func (e *Engine) InviteUser(req *InviteRequest) (*domain.UserInvitation, error) {
	if _, err := e.sqlClient.FindUser(domain.User{ID: req.FromUserID}); err != nil {
		return nil, err
	}
	user, err := e.sqlClient.FindUser(domain.User{
		Email:  req.Email,
		Status: domain.UserStatusTemp,
	})
	if err == nil && user != nil {
		invitation, err := e.sqlClient.FindInvitation(domain.UserInvitation{
			FromUserID: req.FromUserID,
			ToUserID:   user.ID,
			Status:     domain.InvitationStatusSent,
		})
		if err != nil {
			return nil, err
		}
		params := &domain.UserInvitation{Hash: domain.GenerateSlug(64)}
		err = e.sqlClient.UpdateInvitation(invitation, params)
		return invitation, err
	}

	invitation := &domain.UserInvitation{
		Status:     domain.InvitationStatusSent,
		FromUserID: req.FromUserID,
		ToUser: &domain.User{
			Status: domain.UserStatusTemp,
			Email:  req.Email,
		},
	}
	err = e.sqlClient.CreateInvitation(invitation)
	return invitation, err
}
