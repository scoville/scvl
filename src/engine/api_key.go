package engine

import (
	"github.com/scoville/scvl/src/domain"
)

// FindAPIKey returns the api key
func (e *Engine) FindAPIKey(value string) (apiKey *domain.APIKey, err error) {
	apiKey, err = e.sqlClient.FindAPIKey(value)
	return
}
