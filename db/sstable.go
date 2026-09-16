package db

import (
	"fmt"
	"lsm-tree/memtable"
	"os"
	"path/filepath"
	"sort"
)


func Flush(mem *memtable.MemTable, sstable_incr int, path string) error {
	keys := []string{}

	filename := fmt.Sprintf("seg-%d.txt",sstable_incr)
	filePath := filepath.Join(path,filename)
	file,err := os.OpenFile(filePath,os.O_CREATE | os.O_WRONLY | os.O_TRUNC,0644)

	if err != nil {
		return err
	}
	
	defer file.Close()

	fmt.Println(file)
	for key := range mem.Entries(){
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _,key := range keys {
		value := mem.Entries()[key]

		_,err := file.Write([]byte(key + "|" + string(value) + "\n"))

		if err != nil {
			return err
		}
	}

	fmt.Println(keys)

	return nil
}