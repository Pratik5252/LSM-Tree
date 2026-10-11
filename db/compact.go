package db

import (
	"fmt"
	"lsm-tree/memtable"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func (db *DB) Compact(path string) error {
	merged := make(map[string]memtable.Entry)
	ssTableIndex := db.nextSSTableID

	files, err := filepath.Glob(filepath.Join(path,"seg-*.txt"))

	if err != nil {
		return err
	}
	
	if len(files) == 0 {
		return nil
	}

	sort.Slice(files,func(i, j int) bool {
		return getSSTableID(files[i]) < getSSTableID(files[j])
	})

	for _, file := range files {
		records, err := ReadAll(file)

		if err != nil {
			return err
		}
		for key, value := range records {
			merged[key] = value
		}
		
	}

	err = db.mergeSSTable(path, merged)

	if err != nil {
		return err
	}
	
	for _, file := range files {
		if getSSTableID(file) < ssTableIndex{
			if err := os.Remove(file); err != nil{
				return err
			}
		}
	}
	db.nextSSTableID++
	
	return nil
}

func getSSTableID(file string) int {
	filename := filepath.Base(file)
	prefix := strings.TrimPrefix(filename,"seg-")
	suffix := strings.TrimSuffix(prefix,".txt")
	id, _ := strconv.Atoi(suffix)

	return id
}

func (db *DB)mergeSSTable(path string, merged map[string]memtable.Entry) error{

	keys := []string{}

	filename := fmt.Sprintf("seg-%d.txt", db.nextSSTableID)
	filepath := filepath.Join(path,filename)
	file, err := os.OpenFile(filepath,os.O_CREATE | os.O_WRONLY | os.O_TRUNC,0644)

	if err != nil {
		return err
	}
	

	for key := range merged{
		if merged[key].Tombstone == true {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _,key := range keys {
		entry := merged[key]

		_,err := file.Write([]byte(key + "|" + string(entry.Value) + "|" + strconv.FormatBool(entry.Tombstone) + "\n"))

		if err != nil {
			file.Close()
			return err
		}
	}

	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}
	return nil
}