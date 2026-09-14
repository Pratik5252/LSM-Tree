package db

import (
	"fmt"
	"lsm-tree/memtable"
	"sort"
)

func Flush(mem *memtable.MemTable) error {

	
	keys := []string{}
	
	for key := range mem.Entries(){
		keys = append(keys, key)
	}

	sort.Strings(keys)

	fmt.Println(keys)

	return nil
}