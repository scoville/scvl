package engine

import "github.com/scoville/scvl/src/domain"

// SQLClient is the interface for sql client
type SQLClient interface {
	FindUser(uint) (*domain.User, error)
}
