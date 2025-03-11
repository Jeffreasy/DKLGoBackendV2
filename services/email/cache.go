package email

import (
	"dklautomationgo/models"
	"sync"
	"time"
)

// AccountCache holds the cache for one account
type AccountCache struct {
	emails     []*models.Email
	lastFetch  time.Time
	cacheMutex sync.RWMutex
}

func NewAccountCache() *AccountCache {
	return &AccountCache{
		emails:     make([]*models.Email, 0),
		lastFetch:  time.Time{},
		cacheMutex: sync.RWMutex{},
	}
}
