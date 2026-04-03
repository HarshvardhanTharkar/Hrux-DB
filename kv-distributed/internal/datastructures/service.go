package datastructures

import (
	"errors"
	"sort"
	"sync"
)

type Set map[string]struct{}
type SortedList []string
type Map map[string]interface{}
type Queue []string
type Stack []string

type DataStructuresService struct {
	mu     sync.RWMutex
	Sets   map[string]Set
	Lists  map[string]SortedList
	Maps   map[string]Map
	Queues map[string]Queue
	Stacks map[string]Stack
}

func NewDataStructuresService() *DataStructuresService {
	return &DataStructuresService{
		Sets:   make(map[string]Set),
		Lists:  make(map[string]SortedList),
		Maps:   make(map[string]Map),
		Queues: make(map[string]Queue),
		Stacks: make(map[string]Stack),
	}
}

func (ds *DataStructuresService) SetAdd(setName, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.Sets[setName]; !ok {
		ds.Sets[setName] = make(Set)
	}
	ds.Sets[setName][value] = struct{}{}
}

func (ds *DataStructuresService) SetRemove(setName, value string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if s, ok := ds.Sets[setName]; ok {
		if _, exists := s[value]; exists {
			delete(s, value)
			return nil
		}
		return errors.New("value not found in set")
	}
	return errors.New("set not found")
}

func (ds *DataStructuresService) SetList(setName string) ([]string, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if s, ok := ds.Sets[setName]; ok {
		var list []string
		for v := range s {
			list = append(list, v)
		}
		return list, nil
	}
	return nil, errors.New("set not found")
}

func (ds *DataStructuresService) SortedListAdd(listName, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Lists[listName] = append(ds.Lists[listName], value)
	sort.Strings(ds.Lists[listName])
}

func (ds *DataStructuresService) SortedListGet(listName string) ([]string, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if l, ok := ds.Lists[listName]; ok {
		return l, nil
	}
	return nil, errors.New("list not found")
}

func (ds *DataStructuresService) MapPut(mapName, key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, ok := ds.Maps[mapName]; !ok {
		ds.Maps[mapName] = make(Map)
	}
	ds.Maps[mapName][key] = value
}

func (ds *DataStructuresService) MapGet(mapName, key string) (interface{}, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if m, ok := ds.Maps[mapName]; ok {
		if val, ok := m[key]; ok {
			return val, nil
		}
		return nil, errors.New("key not found in map")
	}
	return nil, errors.New("map not found")
}

func (ds *DataStructuresService) QueuePush(name, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Queues[name] = append(ds.Queues[name], value)
}

func (ds *DataStructuresService) QueuePop(name string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	q := ds.Queues[name]
	if len(q) == 0 {
		return "", errors.New("queue empty")
	}
	val := q[0]
	ds.Queues[name] = q[1:]
	return val, nil
}

func (ds *DataStructuresService) QueuePeek(name string) (string, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	q := ds.Queues[name]
	if len(q) == 0 {
		return "", errors.New("queue empty")
	}
	return q[0], nil
}

func (ds *DataStructuresService) StackPush(name, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Stacks[name] = append(ds.Stacks[name], value)
}

func (ds *DataStructuresService) StackPop(name string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	s := ds.Stacks[name]
	if len(s) == 0 {
		return "", errors.New("stack empty")
	}
	val := s[len(s)-1]
	ds.Stacks[name] = s[:len(s)-1]
	return val, nil
}

func (ds *DataStructuresService) StackPeek(name string) (string, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	s := ds.Stacks[name]
	if len(s) == 0 {
		return "", errors.New("stack empty")
	}
	return s[len(s)-1], nil
}
