package wal

import (
	"bufio"
	"lsm-tree/record"
	"os"
)

type WAL struct {
	file *os.File
	error
}

func Open(path string) (*WAL,error){
	wal, err := os.OpenFile(path,os.O_CREATE | os.O_RDWR | os.O_APPEND,0644)

	return &WAL{
		file: wal,
	}, err
}

func (wal *WAL) append(line []byte) error{
	_, err := wal.file.Write(line)

	if err != nil {
		return err
	}

	return nil
}

// implement PUT that add a log to wal.log
func (wal *WAL) Put(key string, value []byte) error{
	line := []byte("PUT | " + key + " | " + string(value) + "\n")
	return wal.append(line)
}

// Implement DELETE taht add a delete record log in wal.log
func (wal * WAL) Delete(key string) error{
	line := []byte("DELETE | " + key + "\n")
	return wal.append(line)
}

func(wal *WAL) Replay(fn func(*record.Record) error) error{
	scanner := bufio.NewScanner(wal.file)

	for scanner.Scan(){
		line := scanner.Bytes()
		rec,err := record.ParseWal(line)

		if err != nil{
			return err
		}

		if err := fn(rec); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil{
		return err
	}

	return nil
}