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

	db.Put("apple", []byte("red"))
	db.Delete("apple")
	db.Put("apple", []byte("green"))
	db.Put("button", []byte("blue"))
	db.Put("orange", []byte("red"))
	db.Put("lamp", []byte("yellow"))
	db.Put("curtain", []byte("purple"))
	db.Put("dog", []byte("white"))
	db.Delete("curtain")
	db.Put("dog", []byte("black"))
	db.Put("lantern", []byte("green"))
	db.Put("mouse", []byte("grey"))

	db.Compact("./data")
	entry, ok, err := db.Get("curtain")

	if err != nil {
		fmt.Println(err)
	}

	if !ok {
		fmt.Println("Key does not exist")
	} else if entry.Tombstone {
		fmt.Println("Key was deleted")
	} else {
		fmt.Println(string(entry.Value))
	}
	
}