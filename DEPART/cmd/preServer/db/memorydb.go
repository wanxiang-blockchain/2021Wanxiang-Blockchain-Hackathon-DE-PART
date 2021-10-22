package db

import (
	"preServer/utils"
	"sync"
)

type memoryDB struct {
	lock sync.Mutex
	data map[string]*utils.Record
}

func NewMemoryDB() *memoryDB {
	return &memoryDB{
		data: make(map[string]*utils.Record, 1024),
	}
}

func (m memoryDB) Put(record *utils.Record) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[record.ID] = record
}

func (m memoryDB) Get(id string) *utils.Record {
	m.lock.Lock()
	defer m.lock.Unlock()
	if rec, ok := m.data[id]; ok {
		return rec
	}
	return nil
}
