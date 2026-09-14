package db

import (
	"lsm-tree/memtable"
	"lsm-tree/record"
	"lsm-tree/wal"
	"path/filepath"
)

// Open opens a database at the specified path and returns a DB instance
func Open(path string) (*DB, error) {

	wal, err := wal.Open(filepath.Join(path,"wal.log"))

	if err != nil {
		return nil, err
	}
	
	memtable := memtable.New()


	wal.Replay(func(rec *record.Record) error {
		switch rec.Operation {
			case record.OpPut:
				memtable.Put(rec.Key,rec.Value)
			case record.OpDelete:
				memtable.Delete(rec.Key)
		}
		return nil
	})

	
	return &DB{
		memtable: memtable,
		wal: wal,
	}, nil
}