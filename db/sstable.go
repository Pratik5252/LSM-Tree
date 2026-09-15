package db

import (
	"fmt"
	"lsm-tree/memtable"
	"os"
	"sort"
)


func Flush(mem *memtable.MemTable) error {
	keys := []string{}

	file,err := os.OpenFile("./data/seg.txt",os.O_CREATE | os.O_WRONLY | os.O_TRUNC,0644)
	if err != nil{
		return err
	}

	fmt.Println(file)
	for key := range mem.Entries(){
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _,key := range keys {
		value := mem.Entries()[key]

		file.Write([]byte(key + "|" + string(value) + "\n"))
	}

	fmt.Println(keys)

	return nil
}