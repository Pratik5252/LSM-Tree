package db

import (
	"fmt"
	"lsm-tree/memtable"
	"lsm-tree/wal"
	"path/filepath"
	"sort"
)

type DB struct {
	memtable *memtable.MemTable
	wal *wal.WAL
	nextSSTableID int
	path string
}

// Put adds key-value pair to the database
func (db *DB) Put(key string, value []byte){

	var MAX_MEMTABLE_SIZE = 1

	err := db.wal.Put(key,value);
	if err != nil {
		panic(err)
	}
	db.memtable.Put(key,value)

	if db.memtable.Count() > MAX_MEMTABLE_SIZE {
		err := Flush(db.memtable,db.nextSSTableID,db.path)

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

	fmt.Println("Get", key, entry, ok)
	if ok {
		return entry,true, nil
	}

	// If not found in memtable, search in SSTables
	files, err := filepath.Glob(filepath.Join(db.path, "seg-*.txt"))
	
	if err != nil {
		return memtable.Entry{}, false, err
	}

	sort.Slice(files,func(i, j int) bool {
		return getSSTableID(files[i]) < getSSTableID(files[j])
	})
	
	for _,file := range files {
		entry,ok,err := Read(file,key)

		if err != nil {
			return memtable.Entry{}, false, err
		}

		if entry.Tombstone {
			return entry, true, nil
		}
		
		if ok {
			return entry,true, nil
		}
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