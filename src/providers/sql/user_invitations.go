package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/scoville/scvl/src/domain"
)

const tblUserInvitations = "user_invitations"

func (c *client) FindInvitation(hash string) (*domain.UserInvitation, error) {
	invitation := &domain.UserInvitation{}
	err := c.db.Table(tblUserInvitations).
		Preload("ToUser", func(db *gorm.DB) *gorm.DB {
			return db
		}).First(invitation, "hash = ?", hash).Error
	return invitation, err
}

func (c *client) UpdateInvitation(invitation, params *domain.UserInvitation) (*domain.UserInvitation, error) {
	err := c.db.Model(invitation).Updates(params).Error
	return invitation, err
}

func (c *client) CreateInvitation(params *domain.UserInvitation) (*domain.UserInvitation, error) {
	err := c.db.Create(params).Error
	return params, err
}
