package sql

import (
	"github.com/scoville/scvl/src/domain"
)

const tblAPIKeys = "api_keys"

func (c *client) FindAPIKey(value string) (apiKey *domain.APIKey, err error) {
	apiKey = &domain.APIKey{}
	err = c.db.Table(tblAPIKeys).
		Where("status = ?", domain.APIKeyStatusActive).
		First(apiKey, "value = ?", value).Error
	return
}
