package http

import (
	"github.com/graphicweave/injun/database"
	"context"
)

func AddToWishlist(userId, eventId string) error {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	if err != nil {
		return err
	}

	return db.UpdateWishlist(userId, eventId, true)
}

func RemoveToWishlist(userId, eventId string) error {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	if err != nil {
		return err
	}

	return db.UpdateWishlist(userId, eventId, false)
}

func GetWislistEvents(userId string) ([]map[string]interface{}, error) {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	if err != nil {
		return nil, err
	}

	return db.GetWishlist(userId)
}

func getStats(userId string) (map[string]float64, error) {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	if err != nil {
		return nil, err
	}

	return db.GetStats(userId)
}
