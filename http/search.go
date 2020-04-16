package http

import (
	"net/http"
	"github.com/graphicweave/injun/elastic"
	"github.com/valyala/fasthttp"
	log "github.com/sirupsen/logrus"
	"context"
	"encoding/json"
	"github.com/graphicweave/injun/service"
	"fmt"
	"reflect"
)

func SearchHandler(ctx *fasthttp.RequestCtx) {

	uri := ctx.URI()
	var events []string

	if string(uri.Path()) == "/search" && uri.QueryArgs().Has("q") {

		searchTerm := string(uri.QueryArgs().Peek("q"))
        userId :=  string(uri.QueryArgs().Peek("id"))
        fmt.Println("st",searchTerm)
        fmt.Println("Uid",userId)

		elasticSearch, err := elastic.NewElasticSearch(context.Background())
		if err != nil {
			log.Error("failed to get elastic search client")
			log.Errorln("Error ", err)
			ctx.Error("internal server error", http.StatusInternalServerError)
			return
		}
		if userId != "" {
			arangoDb, err := service.NewArangoDB(context.Background())
			if err != nil {
				if err != nil {
					log.Error("failed to get ArangoDB connection")
					log.Error("Error: ", err)
					ctx.Error("internal server error", http.StatusInternalServerError)
					return
				}
			}
			events, err = arangoDb.GetWishlisted(userId)
			if err != nil {
				log.Error("failed to get wishlist of the user")
				log.Errorln("Error ", err)
				ctx.Error("internal server error", http.StatusInternalServerError)
				return
			}
		}

		results, err := elasticSearch.Search(searchTerm, "#", "#","#", events, true)
		if err != nil {
			log.Error("failed to search for term: ", searchTerm)
			log.Errorln("Error ", err)
			ctx.Error("internal server error", http.StatusInternalServerError)
			return
		}

		if err == elastic.ErrEsNoResults {
			log.Error("no search results found: ", searchTerm)
			ctx.WriteString(err.Error())
			return
		}

		json.NewEncoder(ctx).Encode(results)
		return

	} else if ctx.IsPost() && string(uri.Path()) == "/search" {

		var args SearchTerm
		json.Unmarshal(ctx.Request.Body(), &args)
		if args.UserId != " " {
			arangoDb, err := service.NewArangoDB(context.Background())
			if err != nil {
				log.Error("failed to get arangodb client")
				log.Errorln("Error ", err)
			}
			events, err = arangoDb.GetWishlisted(args.UserId)
			if err != nil {
				log.Error("failed to get wishlist events")
				log.Errorln("Error ", err)
			}
		}
		elasticSearch, err := elastic.NewElasticSearch(context.Background())
		if err != nil {
			log.Error("failed to get elastic search client")
			log.Errorln("Error ", err)
			ctx.Error("internal server error", http.StatusInternalServerError)
			return
		}

		results, err := elasticSearch.Search(args.SearchTerm, args.Country, args.State, args.City, events,false)

		if err != nil && err != elastic.ErrEsNoResults {
			log.Error("failed to search for term: ", args.SearchTerm)
			log.Errorln("Error ", err)
			ctx.Error("internal server error", http.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(http.StatusOK)
		json.NewEncoder(ctx).Encode(results)
	} else {
		ctx.Error(NotFoundError.Error(), http.StatusNotFound)
	}
}

func NearbyEventsHandler(ctx *fasthttp.RequestCtx) {

	uri := ctx.URI()

	lat, _ := uri.QueryArgs().GetUfloat("lat")
	lon, _ := uri.QueryArgs().GetUfloat("lon")

	fmt.Println("lat", lat+lon)

	elasticSearch, err := elastic.NewElasticSearch(context.Background())

	if err != nil {
		log.Error("failed to get elastic search client")
		log.Errorln("Error ", err)
		ctx.Error("internal server error", http.StatusInternalServerError)
		return
	}

	results, err := elasticSearch.NearbyEvent(lat, lon)

	if err != nil || reflect.ValueOf(results).Len() == 0  {
		log.Error("failed to search nearby events", )
		log.Errorln("Error ", err)

		arangoDb, err := service.NewArangoDB(context.Background())

		if err != nil {
			log.Error("failed to get ArangoDB connection")
			log.Error("Error: ", err)

		}
		events, err := arangoDb.GetAlternateEvents()

		if err != nil {
			log.Error("failed to get Alternate Event")
			log.Error("Error: ", err)
		}
		ctx.SetStatusCode(http.StatusOK)
		json.NewEncoder(ctx).Encode(map[string]interface{}{"nearby": false , "events" : events})
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	json.NewEncoder(ctx).Encode(map[string]interface{}{"nearby": true , "events" : results})
}

type SearchTerm struct {
	SearchTerm   string
	State   string
	Country string
	City string
	UserId string
}
