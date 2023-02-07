package tests

import (
	"errors"
	"strings"
	"testing"

	"com.blocopad/blocopad_poc/internal/backend"
	"com.blocopad/blocopad_poc/internal/db"
)

var (
	deleteInvoked bool
	deletedKey    string
)

func TestGetKeyOk(t *testing.T) {
	// Given
	db.GetNote = func(key string) (bool, string, error) {
		return false, "OK", nil
	}

	// When
	data, err := backend.GetKey("key1")

	// Then
	if err != nil {
		t.Fatal("TestGetKeyOk Should not return error")
	}
	if data != "OK" {
		t.Fatal("TestGetKeyOk Invalid return")
	}
}

func TestGetErrorSize(t *testing.T) {
	// Given
	keyString := ""
	keyString2 := strings.Repeat("a", 40)

	// When
	_, err := backend.GetKey(keyString)

	// Then
	if err == nil {
		t.Fatal("TestGetErrorSize Should return error zero length key")
	}

	// When
	_, err = backend.GetKey(keyString2)

	// Then
	if err == nil {
		t.Fatal("TestGetErrorSize Should return error key > 36")
	}
}

func TestGetKeyDbError(t *testing.T) {
	// Given
	db.GetNote = func(key string) (bool, string, error) {
		return false, "OK", errors.New("Error")
	}

	// When
	_, err := backend.GetKey("key1")

	// Then
	if err == nil {
		t.Fatal("TestGetKeyDbError Should return error")
	}
}

func TestSaveKeyOK(t *testing.T) {

	// Given
	db.SaveNote = func(data string, oneTime bool) (string, error) {
		return "123456", nil
	}

	// When
	uuid, err := backend.SaveKey("blablabla", false)

	// Then
	if err != nil {
		t.Fatal("TestSaveKey OK Should not return error")
	}

	if uuid != "123456" {
		t.Fatal("TestSaveKeyOK Invalid return")
	}
}

func TestSaveKeyDbError(t *testing.T) {

	// Given
	db.SaveNote = func(data string, oneTime bool) (string, error) {
		return "123456", errors.New("Error")
	}

	// When
	_, err := backend.SaveKey("blablabla", false)

	// Then
	if err == nil {
		t.Fatal("TestSaveKeyDbError Should return error")
	}
}

func TestSaveInvalidSize(t *testing.T) {

	// Given
	dataZeroLength := ""
	dataTooBig := strings.Repeat("a", 330000)

	// When
	_, err := backend.SaveKey(dataZeroLength, false)

	// Then
	if err == nil {
		t.Fatal("TestSaveInvalidSize it should return an error on zero length")
	}

	// When
	_, err = backend.SaveKey(dataTooBig, false)

	// Then
	if err == nil {
		t.Fatal("TestSaveInvalidSize it should return an error on too big value")
	}

}

func TestGetKeyDeleteOk(t *testing.T) {

	// Given
	deleteInvoked = false
	deletedKey = ""

	db.GetNote = func(key string) (bool, string, error) {
		return true, "OK", nil
	}

	db.DeleteNote = func(key string) error {
		deleteInvoked = true
		deletedKey = key
		return nil
	}

	// When
	data, err := backend.GetKey("key1")

	// Then
	if err != nil {
		t.Fatal("TestGetKeyDeleteOk Should not return error")
	}
	if data != "OK" {
		t.Fatal("TestGetKeyDeleteOk Invalid return")
	}
	if !deleteInvoked {
		t.Fatal("TestGetKeyDeleteOk should have deleted the key")
	}
	if deletedKey != "key1" {
		t.Fatal("TestGetKeyDeleteOk should have deleted the right key")
	}
}

func TestGetKeyDeleteDbError(t *testing.T) {

	// Given
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("TestGetKeyDeleteDbError The code did not panic")
		}
	}()

	deleteInvoked = false
	deletedKey = ""

	db.GetNote = func(key string) (bool, string, error) {
		return true, "OK", nil
	}

	db.DeleteNote = func(key string) error {
		deleteInvoked = true
		deletedKey = key
		return errors.New("Error")
	}

	// When
	_, err := backend.GetKey("key1")

	// Then
	if err != nil {
		t.Fatal("TestGetKeyDeleteDbError Should not return error")
	}
}
