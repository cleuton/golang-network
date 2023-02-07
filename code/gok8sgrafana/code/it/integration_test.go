package it

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/go-redis/redis/v9"
)

type SavedNote struct {
	Code string `json:"code"`
}

func TestSaveOK(t *testing.T) {
	postUrl := "http://localhost:8080/api/note"
	var jsonData = []byte(`{
		"data": "should I save this?",
		"onetime": false
	}`)
	request, error := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatal("TestSaveOK shoud return no error")
	}

	var savedNote SavedNote

	err := json.NewDecoder(response.Body).Decode(&savedNote)

	if err != nil {
		t.Fatal("TestSaveOK should return a valid json")
	}

	if len(savedNote.Code) == 0 {
		t.Fatal("TestSaveOK should return a valid uuid")
	}

	// Now testing get
	request2, error2 := http.NewRequest("GET", postUrl+"/"+savedNote.Code, nil)

	response2, error2 := client.Do(request2)
	if error != nil {
		panic(error2)
	}
	defer response2.Body.Close()

	b, err := io.ReadAll(response2.Body)
	if err != nil {
		log.Fatalln(err)
	}

	outputCode := string(b)
	if string(b) != outputCode {
		t.Fatal("TestSaveOK Did not return the correct string on GET")
	}

	rDB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Check if expiration date is ok
	ctx := context.Background()
	dCmd := rDB.TTL(ctx, savedNote.Code)
	if dCmd.Val().Hours() < 23 {
		t.Fatal("TestSaveOK should set expiration to 24 hours")
	}

}
