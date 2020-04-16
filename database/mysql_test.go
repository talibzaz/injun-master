package database

import (
	"testing"
	"encoding/json"
	"os"
)

func TestGetPositions(t *testing.T) {
	emails := []interface{}{"john@doe.com", "asmayasrib@gmail.com"}
	positions, err := GetPositions(len(emails), emails)

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", " ")
	e.Encode(positions)
}