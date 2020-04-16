package http

import (
	"github.com/graphicweave/injun/database"
	"context"
	"github.com/graphicweave/injun/elastic"
)

/**
 * Created by  ♅ Salfi Farooq on 23/06/17.
 */
// TODO Refactor.
func PublishEvent(arangoId, elasticId, date string) error {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)
	if err != nil {
		return err
	}

	err  = db.PublishEvent(arangoId, date)
	if err != nil {
		return err
	}
	elasticSearch, err := elastic.NewElasticSearch(context.Background())
	if err != nil {
		return err
	}

	err = elasticSearch.PublishEvent(elasticId, date)

	if err != nil {
		return err
	}
	return nil
}
