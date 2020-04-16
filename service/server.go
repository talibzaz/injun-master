package service

import (
	"github.com/graphicweave/injun/elastic"
	"context"
	"github.com/graphicweave/injun/database"
)

type Server struct {
	*elastic.ElasticSearch
	*database.ArangoDB
}

// NewServer creates ElasticSearch, ArangoDB clients
func NewServer() (Server, error) {

	server := Server{}
	var err error

	ctx := context.Background()

	// init ElasticSearch
	server.ElasticSearch, err = elastic.NewElasticSearch(ctx)

	// init ArangoDB
	server.ArangoDB, err = database.NewArangoDB(ctx)
	return server, err
}

// NewElasticSearch creates only ElasticSearch client
func NewElasticSearch(ctx context.Context) (*elastic.ElasticSearch, error) {
	return elastic.NewElasticSearch(ctx)
}

// NewArangoDB creates only NewArangoDB client
func NewArangoDB(ctx context.Context) (*database.ArangoDB, error) {
	return database.NewArangoDB(ctx)
}
