package storage

import (
	"sync"
)

type MapKVStorage struct {
	sync.RWMutex
	BasicKVStorage
	values map[string]string
}

func NewMapStorage() *MapKVStorage {
	return &MapKVStorage{
		values: map[string]string{},
	}
}

func (ms *MapKVStorage) Add(key string, value string) (string, error) {
	ms.Lock()
	defer ms.Unlock()
	_, ok := ms.values[key]
	if ok {
		return key, ErrKeyExists
	}
	ms.values[key] = value
	return key, nil
}

func (ms *MapKVStorage) Get(key string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	value, ok := ms.values[key]
	if !ok {
		return key, ErrKeyNotFound
	}
	return value, nil
}

func (ms *MapKVStorage) Del(key string) error {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.values, key)
	return nil
}
