package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/scoville/scvl/src/domain"
)

const tblUserInvitations = "user_invitations"

func (c *client) FindInvitation(cond domain.UserInvitation) (*domain.UserInvitation, error) {
	invitation := &domain.UserInvitation{}
	err := c.db.Table(tblUserInvitations).
		Preload("ToUser", func(db *gorm.DB) *gorm.DB {
			return db
		}).First(invitation, cond).Error
	return invitation, err
}

func (c *client) UpdateInvitation(invitation, params *domain.UserInvitation) (err error) {
	err = c.db.Table(tblUserInvitations).
		Model(invitation).Updates(params).Error
	return err
}

func (c *client) CreateInvitation(params *domain.UserInvitation) (err error) {
	err = c.db.Create(params).Error
	return
}
