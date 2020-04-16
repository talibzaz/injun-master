package event

import (
	"testing"
	"github.com/spf13/viper"
	"github.com/graphicweave/injun/service"
	"context"
	"encoding/json"
	"os"
	"github.com/arangodb/go-driver"
	"fmt"
)

func TestManageEventById(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	arangoDb, err := service.NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	e, err := arangoDb.ManageEventById("bfaofkqbil1g00csa780")

	if err != nil && driver.IsNoMoreDocuments(err) {
		fmt.Println(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(e)
}

func TestEventService_UpdateEventById(t *testing.T) {

}
