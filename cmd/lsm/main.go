package main

import (
	"lsm-tree/db"
	"fmt"
)

func main() {
	db, err := db.Open("./data") 
	if err != nil {
		panic(err)
	}

	db.Put("apple",[]byte("red"))
	db.Put("rose",[]byte("red"))
	db.Put("candle",[]byte("blue"))

	entry,ok,err := db.Get("apple")


	db.Delete("apple")

	if err != nil {
		fmt.Println(err)
	}

	if !ok {
    fmt.Println("Key does not exsist")
	} else if entry.Tombstone {
    fmt.Println("Key was deleted")
	} else {
		fmt.Println(string(entry.Value))
	}
}