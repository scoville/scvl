package domain

import (
	"fmt"
	"time"
)

// UserInvitation is the struct for user_invitation.
type UserInvitation struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Hash       string     `json:"hash" gorm:"unique_index; not null"`
	Status     string     `json:"-" gorm:"not null" valid:"required,in(sent|used|deleted)"`
	FromUserID uint       `json:"-" gorm:"type:integer REFERENCES users(id) ON DELETE CASCADE; not null" valid:"required"`
	ToUserID   uint       `json:"to_user_id" gorm:"type:integer REFERENCES users(id) ON DELETE CASCADE"`

	ToUser *User `json:"to_user,omitempty" gorm:"association_autupdate:false;association_autcreate:false"`
}

// Invitation statuses
const (
	InvitationStatusSent    = "sent"
	InvitationStatusUsed    = "used"
	InvitationStatusDeleted = "deleted"
)

// BeforeCreate generates a unique hash for the invitation.
func (i *UserInvitation) BeforeCreate() error {
	i.Hash = GenerateSlug(64)
	i.Status = InvitationStatusSent
	i.ToUser.SetPassword(GenerateSlug(12))
	return nil
}

// Valid returns Error if the invitation is not valid
func (i *UserInvitation) Valid() error {
	if i.Status == InvitationStatusUsed {
		return fmt.Errorf("the invitation is already used")
	}
	if time.Now().Sub(i.CreatedAt) > time.Hour*24 {
		return fmt.Errorf("the invitation is expired")
	}
	return nil
}
