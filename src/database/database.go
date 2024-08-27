package database

import (
	"fmt"
)


var db map[string]string = make(map[string]string)

var ErrNoKey = fmt.Errorf("no such key")


func SetKey(key, val string) {
	db[key] = val
}


func GetKey(key string) (string, error) {
	val, ok := db[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrNoKey, key)
	} else {
		return val, nil
	}
}
