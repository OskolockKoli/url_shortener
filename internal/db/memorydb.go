package db

import (
	"sync"

	"github.com/OskolockKoli/url_shortener/internal/models"
)

type MemoryDB struct {
	mu      sync.RWMutex
	storage map[string]models.Link
}

func (m *MemoryDB) Save(link models.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.storage[link.ShortID]; exists {
		return ErrDuplicateKey
	}

	m.storage[link.ShortID] = link
	return nil
}

func (m *MemoryDB) GetByShortID(shortID string) (models.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	link, found := m.storage[shortID]
	if !found {
		return models.Link{}, ErrNotFound
	}

	return link, nil
}

func (m *MemoryDB) Init() {
	m.storage = make(map[string]models.Link)
}
