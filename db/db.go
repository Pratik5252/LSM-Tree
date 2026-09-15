package db

import (
	"lsm-tree/memtable"
	"lsm-tree/wal"
)

type DB struct {
	memtable *memtable.MemTable
	wal *wal.WAL
}

// Put adds key-value pair to the database
func (db *DB) Put(key string, value []byte){

	err := db.wal.Put(key,value);
	if err != nil {
		panic(err)
	}
	db.memtable.Put(key,value)

	if db.memtable.Count() > 5 {
		Flush(db.memtable)
	}
}

// Get retrieves value by key and returns value in []byte and a bool as status
func (db *DB) Get(key string)([]byte,bool){
	return db.memtable.Get(key)
}

// Delete removes a key-value pair from the database
func (db *DB) Delete(key string){

	err := db.wal.Delete(key)

	if err != nil {
		panic(err)
	}
	db.memtable.Delete(key)	
}