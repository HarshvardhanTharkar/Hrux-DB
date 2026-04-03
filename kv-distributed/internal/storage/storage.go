package storage

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
	"time"
)

type Storage struct {
	mu       sync.RWMutex
	Buckets  map[string]map[string][]byte
	TTL      map[string]map[string]TTLInfo
	filename string
}

type TTLInfo struct {
	Expiry int64
}

func NewStorage() *Storage {
	return &Storage{
		Buckets: make(map[string]map[string][]byte),
		TTL:     make(map[string]map[string]TTLInfo),
	}
}

func (s *Storage) SetPersistenceFile(filename string) {
	s.filename = filename
}

func (s *Storage) Put(bucket, key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Buckets[bucket]; !ok {
		s.Buckets[bucket] = make(map[string][]byte)
	}
	s.Buckets[bucket][key] = value
}

func (s *Storage) PutWithTTL(bucket, key string, value []byte, ttlSeconds int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Buckets[bucket]; !ok {
		s.Buckets[bucket] = make(map[string][]byte)
	}
	s.Buckets[bucket][key] = value

	if ttlSeconds > 0 {
		if _, ok := s.TTL[bucket]; !ok {
			s.TTL[bucket] = make(map[string]TTLInfo)
		}
		s.TTL[bucket][key] = TTLInfo{Expiry: time.Now().Unix() + ttlSeconds}
	}
}

func (s *Storage) Get(bucket, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.checkTTL(bucket, key) {
		return nil, errors.New("key expired")
	}

	if b, ok := s.Buckets[bucket]; ok {
		if val, ok := b[key]; ok {
			return val, nil
		}
	}
	return nil, errors.New("key not found")
}

func (s *Storage) Delete(bucket, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if b, ok := s.Buckets[bucket]; ok {
		if _, exists := b[key]; exists {
			delete(b, key)
			if ttlBucket, ok := s.TTL[bucket]; ok {
				delete(ttlBucket, key)
			}
			return nil
		}
		return errors.New("key not found")
	}
	return errors.New("bucket not found")
}

func (s *Storage) Update(bucket, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if b, ok := s.Buckets[bucket]; ok {
		if _, exists := b[key]; exists {
			b[key] = value
			return nil
		}
		return errors.New("key not found")
	}
	return errors.New("bucket not found")
}

func (s *Storage) List(bucket string) (map[string][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if b, ok := s.Buckets[bucket]; ok {
		result := make(map[string][]byte)
		for k, v := range b {
			if s.checkTTL(bucket, k) {
				result[k] = v
			}
		}
		return result, nil
	}
	return nil, errors.New("bucket not found")
}

func (s *Storage) checkTTL(bucket, key string) bool {
	if bucketTTL, ok := s.TTL[bucket]; ok {
		if t, ok := bucketTTL[key]; ok {
			if time.Now().Unix() > t.Expiry {
				delete(s.Buckets[bucket], key)
				delete(bucketTTL, key)
				return false
			}
		}
	}
	return true
}

func (s *Storage) SaveToFile() error {
	if s.filename == "" {
		return errors.New("no filename set for persistence")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(s)
}

func (s *Storage) LoadFromFile() error {
	if s.filename == "" {
		return errors.New("no filename set for persistence")
	}

	file, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(s)
}
