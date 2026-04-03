package indexing

import (
	"sort"
	"strings"
	"sync"
)

type Indexer struct {
	mu sync.RWMutex
}

func NewIndexer() *Indexer {
	return &Indexer{}
}

func (idx *Indexer) ListKeysByPrefix(data map[string][]byte, prefix string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var keys []string
	for k := range data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (idx *Indexer) ListKeysInRange(data map[string][]byte, start, end string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var keys []string
	for k := range data {
		if k >= start && k <= end {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}
