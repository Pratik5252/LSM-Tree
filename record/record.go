package record

import (
"strings"
"errors"
)

type Operation string

const (
	OpPut Operation = "PUT"
	OpDelete Operation = "DELETE"
)

type Record struct {
	Operation Operation
	Key       string
	Value     []byte
}

// This parse the 
func ParseWal(line []byte) (*Record,error) {
	parts := strings.Split(string(line),"|")

	operator := strings.TrimSpace(parts[0])
	var key string
	var value []byte
	switch operator {
		case "PUT":
			if(len(parts) != 3){
				return nil, errors.New("Not enough parts for PUT operation")
			}
			key = strings.TrimSpace(parts[1])
			value =  []byte(strings.TrimSpace(parts[2]))
		case "DELETE":
			if(len(parts) != 2){
				return nil, errors.New("Not enough parts for DELETE operation")
			}
			key = strings.TrimSpace(parts[1])
		default: 
			return nil, errors.New("Not a valid operator")
	}

	return &Record{
		Operation: Operation(operator),
		Key:       key,
		Value:     value,
	},nil
}