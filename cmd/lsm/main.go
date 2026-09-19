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
	db.Put("ball",[]byte("red"))
	db.Put("light",[]byte("red"))
	db.Put("pen",[]byte("red"))
	db.Put("reel",[]byte("red"))
	db.Put("lamp",[]byte("red"))
	db.Put("bird",[]byte("red"))
	db.Put("fly",[]byte("red"))
	db.Put("aninal",[]byte("red"))
	db.Put("person",[]byte("red"))
	db.Put("blood",[]byte("red"))


	value, ok, err := db.Get("apple")

	if err != nil {
		fmt.Println(err)
	}

	if ok {
    fmt.Println(string(value))
	} else {
    fmt.Println("not found")
	}
}