package db

import (
	"bufio"
	"fmt"
	"lsm-tree/memtable"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

	for key := range mem.Entries(){
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _,key := range keys {
		entry := mem.Entries()[key]

		_,err := file.Write([]byte(key + "|" + string(entry.Value) + "|" + strconv.FormatBool(entry.Tombstone) + "\n"))

		if err != nil {
			return err
		}
	}

	return nil
}

// Read data from segment files
func Read(path string,key string) (memtable.Entry,bool,error){

	file,err := os.Open(path)

	if err != nil {
		return memtable.Entry{},false,err
	}

	defer file.Close()
	
	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		line := scanner.Bytes()

		parts := strings.Split(string(line),"|")

		if len(parts) != 3 {
			continue
		}

		if parts[0] == key{

			if parts[2] == "true" {
				return memtable.Entry{Tombstone: true}, true,nil
			}
			
			fmt.Printf("Value %s",parts[1])
			return memtable.Entry{
				Value: []byte(parts[1]),
				Tombstone: false,
			}, true,nil
		}
	}

	if err := scanner.Err(); err != nil {
		return memtable.Entry{},false,err
	}
	return memtable.Entry{},false,nil
}