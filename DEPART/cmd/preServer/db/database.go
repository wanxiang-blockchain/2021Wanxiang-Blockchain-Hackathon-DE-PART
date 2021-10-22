package db

import "preServer/utils"

var preDB = NewDataBase()

type Database interface {
	Put(record *utils.Record)

	Get(id string) *utils.Record
}

func NewDataBase() Database {
	return NewMemoryDB()
}

func Put(record *utils.Record) {
	preDB.Put(record)
}

func Get(id string) *utils.Record {
	return preDB.Get(id)
}
