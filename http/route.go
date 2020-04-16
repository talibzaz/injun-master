package http

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"reflect"
	"fmt"
)

type AnyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func AddToWishlistHandler(ctx *fasthttp.RequestCtx) {

	postArgs := ctx.Request.PostArgs()

	response := AnyResponse{}
	if !postArgs.Has("user_id") {
		response.Status = "ERR"
		response.Message = "'user_id' missing"
	} else if !postArgs.Has("event_id") {
		response.Status = "ERR"
		response.Message = "'event_id' missing"
	}

	if response != (AnyResponse{}) {
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
			return
		}
	}

	userId := string(postArgs.Peek("user_id"))
	eventId := string(postArgs.Peek("event_id"))

	log.Info("add-wishlist revd user_id: " + userId)
	log.Info("add-wishlist revd event_id: " + eventId)

	err := AddToWishlist(userId, eventId)

	if err != nil {
		log.Error("failed to add event to wishlist")
		log.Error("user_id: " + userId + " , event_id: " + eventId)
		log.Error("Error: " + err.Error())

		response.Status = "ERR"
		response.Message = "failed to add event to wishlist. " + err.Error()
	} else {
		response.Status = "OK"
	}
	if err := json.NewEncoder(ctx).Encode(&response); err != nil {
		httpError(ctx, err)
		return
	}
}

func RemoveFromWishlistHandler(ctx *fasthttp.RequestCtx) {
	postArgs := ctx.Request.PostArgs()

	response := AnyResponse{}

	if !postArgs.Has("user_id") {
		response.Status = "ERR"
		response.Message = "'user_id' missing"
	} else if !postArgs.Has("event_id") {
		response.Status = "ERR"
		response.Message = "'event_id' missing"
	}
	if response != (AnyResponse{}) {
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
			return
		}
	}

	userId := string(postArgs.Peek("user_id"))
	eventId := string(postArgs.Peek("event_id"))

	log.Info("remove-wishlist revd user_id: " + userId)
	log.Info("remove-wishlist revd event_id: " + eventId)

	err := RemoveToWishlist(userId, eventId)

	if err != nil {
		log.Error("failed to remove event from wishlist")
		log.Error("user_id: " + userId + " , event_id: " + eventId)
		log.Error("Error: " + err.Error())

		response.Status = "ERR"
		response.Message = "failed to add event to wishlist. " + err.Error()
	} else {
		response.Status = "OK"
	}

	if err := json.NewEncoder(ctx).Encode(struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "OK",
		Message: "Removed",
	}); err != nil {

		httpError(ctx, err)
	}
}

func GetWishlistHandler(ctx *fasthttp.RequestCtx) {

	uri := ctx.URI()

	response := struct {
		Response AnyResponse
		Wishlist []map[string]interface{}
	}{}

	if uri.QueryArgs().Has("user_id") {

		userId := string(uri.QueryArgs().Peek("user_id"))

		resp := AnyResponse{}

		if len(userId) == 0 {
			log.Error("'user_id' missing from path")
			resp.Status = "ERR"
			resp.Message = "'user_id' missing from path"
			response.Response = resp
		} else {
			wishlist, err := GetWislistEvents(userId)
			if err != nil {
				resp.Status = "ERR"
				resp.Message = err.Error()
			} else {
				response.Wishlist = wishlist
				resp.Status = "OK"
			}
			response.Response = resp
		}
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
			return
		}
	} else {
		log.Error("'user_id' missing")
		resp := AnyResponse{}
		resp.Status = "ERR"
		resp.Message = "'user_id' missing"
		response.Response = resp
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
			return
		}
	}
}

func StatsHandler(ctx *fasthttp.RequestCtx) {

	postArgs := ctx.Request.PostArgs()

	response := Statistics{}
	if !postArgs.Has("user_id") {
		response.Status = "ERR"
		response.Message = "'user_id' missing"
	}
	if !reflect.DeepEqual(response, Statistics{}) {
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
		}
		return
	}

	userId := string(postArgs.Peek("user_id"))
	stats, err := getStats(userId)

	if err != nil {
		response.Status = "ERR"
		response.Message = "failed to get stats " + err.Error()

		log.Error("failed to get stats for user_id: ", userId)
		log.Error("Error: ", err.Error())
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
			return
		}
		return
	}
	response.Status = "OK"
	response.Stats = stats

	if err := json.NewEncoder(ctx).Encode(&response); err != nil {
		httpError(ctx, err)
		return
	}
}

func GetSalesHandler(ctx *fasthttp.RequestCtx) {

	postArgs := ctx.Request.PostArgs()

	response := SalesResponse{}
	if !postArgs.Has("event_id") {
		response.Status = "ERR"
		response.Message = "'event_id' missing"
	}
	if !reflect.DeepEqual(response, SalesResponse{}) {
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
		}
		return
	}

	eventId := string(postArgs.Peek("event_id"))
	sales, err := response.getSales(eventId)

	if err != nil {
		response.Status = "ERR"
		response.Message = "failed to get sales " + err.Error()

		log.Error("failed to get stats for event_id: ", eventId)
		log.Error("Error: ", err.Error())
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
			return
		}
		return
	}
	response.Status = "OK"
	response.Sales = sales

	if err := json.NewEncoder(ctx).Encode(&response); err != nil {
		httpError(ctx, err)
		return
	}
}

func GetFeaturedEventsHandler(ctx *fasthttp.RequestCtx) {

	var userId string

	if ctx.QueryArgs().Has("user_id") {
		userId = string(ctx.QueryArgs().Peek("user_id"))
	}

	var events FeaturedEvents

	if err := json.NewEncoder(ctx).Encode(events.getFeaturedEvents(userId)); err != nil {
		httpError(ctx, err)
		return
	}
}

func GeneralReportByUserIdHandler(ctx *fasthttp.RequestCtx) {

	var userId string

	response := AnyResponse{}

	if ctx.QueryArgs().Has("user_id") {
		userId = string(ctx.QueryArgs().Peek("user_id"))
	} else {
		response.Status = "ERR"
		response.Message = "'user_id' missing"
	}

	if !reflect.DeepEqual(response, AnyResponse{}) {
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
		}
		return
	}

	var reports GeneralReport

	if err := json.NewEncoder(ctx).Encode(reports.getGeneralReport(userId)); err != nil {
		httpError(ctx, err)
		return
	}
}

func GetReportByEventIdHandler(ctx *fasthttp.RequestCtx) {

	var eventId string

	response := AnyResponse{}

	if ctx.QueryArgs().Has("event_id") {
		eventId = string(ctx.QueryArgs().Peek("event_id"))
	} else {
		response.Status = "ERR"
		response.Message = "'event_id' missing"
	}

	if !reflect.DeepEqual(response, AnyResponse{}) {
		if err := json.NewEncoder(ctx).Encode(&response); err != nil {
			httpError(ctx, err)
		}
		return
	}

	var report EventReport

	if err := json.NewEncoder(ctx).Encode(report.getEventReport(eventId)); err != nil {
		httpError(ctx, err)
		return
	}
}

func PublishEventHandler(ctx *fasthttp.RequestCtx) {

	postArgs := ctx.Request.Body()
	fmt.Println("Elastic---", string(postArgs))
	var t map[string]interface{}
	e := json.Unmarshal(postArgs, &t)
	if e != nil {
		fmt.Println(e)
	}

	elasticId := t["elastic_id"].(string)
	arangoId := t["arango_id"].(string)
	date := t["date"].(string)

	err := PublishEvent(arangoId, elasticId, date)
	if err != nil {
		httpError(ctx, err)
		return
	}

	if err := json.NewEncoder(ctx).Encode(struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "OK",
		Message: "Updated",
	}); err != nil {
		httpError(ctx, err)
		return
	}
    log.Info("event published",elasticId)
	return

}
