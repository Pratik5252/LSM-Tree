package db

import (
	"fmt"
	"lsm-tree/memtable"
	"lsm-tree/wal"
)

type DB struct {
	memtable *memtable.MemTable
	wal *wal.WAL
	nextSSTableID int
}

// Put adds key-value pair to the database
func (db *DB) Put(key string, value []byte){

	err := db.wal.Put(key,value);
	if err != nil {
		panic(err)
	}
	db.memtable.Put(key,value)

	if db.memtable.Count() > 5 {
		err := Flush(db.memtable,db.nextSSTableID,"./data")

		if err != nil {
			panic(err)
		}

		db.memtable.Clear()
		db.nextSSTableID++
	}

}

// Get retrieves value by key and returns value in []byte and a bool as status
func (db *DB) Get(key string)([]byte,bool,error){
	value,ok := db.memtable.Get(key)

	sstableID := db.nextSSTableID - 1
	if ok {
		return value,true, nil
	}
	for sstableID > 0 {
		value,ok,err := Read(fmt.Sprintf("./data/seg-%d.txt",sstableID),key)

		if err != nil {
			return nil, false, err
		}
		if ok {
			return value,true, nil
		}
		sstableID--
	}
	return nil,false,nil
}

// Delete removes a key-value pair from the database
func (db *DB) Delete(key string){

	err := db.wal.Delete(key)

	if err != nil {
		panic(err)
	}
	db.memtable.Delete(key)	
}