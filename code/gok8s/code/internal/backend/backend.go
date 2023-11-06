package backend

import (
	"errors"
	"fmt"

	"com.blocopad/blocopad_poc/internal/db"
)

func GetKey(key string) (string, error) {
	if len(key) == 0 || len(key) > 36 {
		return "", errors.New("Key with wrong size")
	}
	oneTime, data, err := db.GetNote(key)
	if err != nil {
		return "", err
	}
	if oneTime {
		if err := db.DeleteNote(key); err != nil {
			panic("Cannot delete onetime note!!!!!")
		}
	}
	return data, nil
}

func SaveKey(data string, oneTime bool) (string, error) {
	byteSize := len([]rune(data))
	if byteSize == 0 || byteSize > (32*1024) {
		return "", errors.New(("Invalid note size"))
	}
	uuidCode, err := db.SaveNote(data, oneTime)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return uuidCode, nil
}
