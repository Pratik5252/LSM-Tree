package db

import (
	"lsm-tree/memtable"
	"lsm-tree/record"
	"lsm-tree/wal"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Open opens a database at the specified path and returns a DB instance
func Open(path string) (*DB, error) {

	err := os.MkdirAll(path,0755)
	if err != nil{
		return nil,err
	}

	wal, err := wal.Open(filepath.Join(path,"wal.log"))

	if err != nil {
		return nil, err
	}
	
	memtable := memtable.New()


	err = wal.Replay(func(rec *record.Record) error {
		switch rec.Operation {
			case record.OpPut:
				memtable.Put(rec.Key,rec.Value)
			case record.OpDelete:
				memtable.Delete(rec.Key)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(path,"seg-*.txt"))

	if err != nil {
		return nil, err
	}

	maxID := 0
	
	for _, file := range files{
		filename := filepath.Base(file)

		prefix := strings.TrimPrefix(filename,"seg-")
		suffix := strings.TrimSuffix(prefix, ".txt")

		id, err := strconv.Atoi(suffix)

		if err != nil {
			return nil, err
		}

		if id > maxID  {
			maxID = id
		}
	}

	currID := maxID + 1
	

	return &DB{
		memtable: memtable,
		wal: wal,
		nextSSTableID: currID,
	}, nil
}