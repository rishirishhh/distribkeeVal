package db_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/rishirishhh/distribkv/db"
)
func TestGetSet(t *testing.T){
	f, err := os.CreateTemp(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp file : %v", err)
	}

	name := f.Name()
	f.Close()
	defer os.Remove(name)

	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create a new Database: %v", err)
	}
	defer closeFunc()

	if err := db.SetKey("party", []byte("Great")); err != nil {
		t.Fatalf("Could not write key : %v", err)
	}

	value, err := db.GetKey("party")
	if err != nil {
		t.Fatalf(`Could not get the key "party": %v`, err )
	}

	if !bytes.Equal(value, []byte("Great")){
		t.Errorf(`Unexpected value for key "party": got %q, want %q`, value, "Greate")
	}
}