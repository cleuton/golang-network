package db

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

var (
	DatabaseUrl      string
	DatabasePassword string
	rDB              *redis.Client
	ctx              = context.Background()
)

var GetDatabase = func() *redis.Client {
	if rDB == nil {
		rDB = redis.NewClient(&redis.Options{
			Addr:     DatabaseUrl,
			Password: DatabasePassword,
			DB:       0,
		})
	}
	return rDB
}

type Note struct {
	Text    string `json:"data"`
	OneTime bool   `json:"onetime"`
}

var GetNote = func(key string) (bool, string, error) {
	db := GetDatabase()
	jsonNote, err := db.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, "", errors.New("not found")
	} else if err != nil {
		// Some other error
		return false, "", err
	}

	var note Note
	err = json.Unmarshal([]byte(jsonNote), &note)
	if err != nil {
		return false, "", err
	}

	return note.OneTime, note.Text, nil
}

var SaveNote = func(data string, oneTime bool) (string, error) {
	stringUuid := (uuid.New()).String()
	db := GetDatabase()
	var note Note
	note.Text = data
	note.OneTime = oneTime
	jsonNote, err := json.Marshal(note)
	if err != nil {
		return "", err
	}
	exp := 24 * time.Hour
	err = db.SetEx(ctx, stringUuid, jsonNote, exp).Err()
	if err != nil {
		return "", err
	}

	return stringUuid, nil
}

var DeleteNote = func(key string) error {
	db := GetDatabase()
	_, err := db.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
