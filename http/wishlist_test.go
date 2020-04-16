package http

import (
	"testing"
	"github.com/spf13/viper"
	"encoding/json"
	"os"
)

func TestAddToWishlist(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	err := AddToWishlist("1288dece-9824-4025-b142-092b3f85dc23", "bd25vvti78ug00f1o0u0")

	if err != nil {
		t.Fatal(err)
	}

}

func TestRemoveToWishlist(t *testing.T) {
	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	err := RemoveToWishlist("1288dece-9824-4025-b142-092b3f85dc23", "bd25vvti78ug00f1o0u0")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetWislistEvents(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	wishlist, err := GetWislistEvents("1288dece-9824-4025-b142-092b3f85dc23")

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&wishlist)
}