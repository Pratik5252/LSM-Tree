package db

import (
	"bufio"
	"fmt"
	"lsm-tree/memtable"
	"os"
	"path/filepath"
	"sort"
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
		value := mem.Entries()[key]

		_,err := file.Write([]byte(key + "|" + string(value) + "\n"))

		if err != nil {
			return err
		}
	}

	return nil
}

// Read data from segment files
func Read(path string,key string) ([]byte,bool,error){

	file,err := os.Open(path)

	if err != nil {
		return nil,false,err
	}

	defer file.Close()
	
	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		line := scanner.Bytes()

		parts := strings.Split(string(line),"|")

		if len(parts) != 2 {
			continue
		}
		if parts[0] == key{
			fmt.Printf("Value %s",parts[1])
			return []byte(parts[1]), true,nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil,false,err
	}
	return nil,false,nil
}