package analytics

import (
	"testing"
	"fmt"
)

func TestGetPageViewByEventID(t *testing.T) {
	res, err := GetPageViewsByEventID("bfd7jhqmc1g00088hag0")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}