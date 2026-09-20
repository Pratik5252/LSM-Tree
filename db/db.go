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
func (db *DB) Get(key string)(memtable.Entry,bool,error){
	entry,ok := db.memtable.Get(key)

	if ok {
		return entry,true, nil
	}

	// If not found in memtable, search in SSTables
	sstableID := db.nextSSTableID - 1
	for sstableID > 0 {
		entry,ok,err := Read(fmt.Sprintf("./data/seg-%d.txt",sstableID),key)

		if err != nil {
			return memtable.Entry{}, false, err
		}

		if entry.Tombstone {
			
		}
		if ok {
			return entry,true, nil
		}
		sstableID--
	}
	return memtable.Entry{},false,nil
}

// Delete removes a key-value pair from the database
func (db *DB) Delete(key string){

	err := db.wal.Delete(key)

	if err != nil {
		panic(err)
	}
	db.memtable.Delete(key)	
}