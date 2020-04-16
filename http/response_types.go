package http

import (
	"github.com/graphicweave/injun/database"
	"context"
	log "github.com/sirupsen/logrus"
)

type Statistics struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Stats   interface{} `json:"stats"`
}

type SalesResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Sales   []database.Sale `json:"sales"`
}

type FeaturedEvents struct {
	Response AnyResponse
	Events   []map[string]interface{}
}

type GeneralReport struct {
	Response AnyResponse
	Report   []map[string]interface{}
}

type EventReport struct {
	Response AnyResponse
	Report   map[string]interface{}
}

func (s SalesResponse) getSales(eventId string) ([]database.Sale, error) {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	if err != nil {
		log.Errorln("failed to create Arango connection")
		log.Errorln("Error:" , err)
		return nil, err
	}

	sales, err := db.GetSales(eventId)

	if err != nil {
		log.Errorln("failed to get sales for event:", eventId)
		log.Errorln("Error:" , err)
		return nil, err
	}

	return sales, nil
}

func (f FeaturedEvents) getFeaturedEvents(userId string) (FeaturedEvents) {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	events := FeaturedEvents{}
	events.Response = AnyResponse{}

	if err != nil {
		log.Errorln("failed to create Arango connection")
		log.Errorln("Error:" , err)
		events.Response.Status = "ERR"
		events.Response.Message = err.Error()
		return events
	}

	resp, err := db.GetFeaturedEvents(userId)

	if err != nil {
		log.Errorln("failed to featured events")
		log.Errorln("userId -", userId , "-")
		log.Errorln("Error:" , err)
		events.Response.Status = "ERR"
		events.Response.Message = err.Error()
		return events
	}

	events.Response.Status = "OK"
	events.Events = resp
	return events
}

func (g GeneralReport) getGeneralReport(userId string) (GeneralReport) {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	report := GeneralReport{}
	report.Response = AnyResponse{}

	if err != nil {
		log.Errorln("failed to create Arango connection")
		log.Errorln("Error:" , err)
		report.Response.Status = "ERR"
		report.Response.Message = err.Error()
		return report
	}

	resp, err := db.GetGeneralReport(userId)

	if err != nil {
		log.Errorln("failed to get report for user:", userId)
		log.Errorln("Error:" , err)
		report.Response.Status = "ERR"
		report.Response.Message = err.Error()
		return report
	}

	report.Response.Status = "OK"
	report.Report = resp
	return report
}

func (e EventReport) getEventReport(eventId string) (EventReport) {

	ctx := context.Background()
	db, err := database.NewArangoDB(ctx)

	report := EventReport{}
	report.Response = AnyResponse{}

	if err != nil {
		log.Errorln("failed to create Arango connection")
		log.Errorln("Error:" , err)
		report.Response.Status = "ERR"
		report.Response.Message = err.Error()
		return report
	}

	resp, err := db.GetReportByEventId(eventId)

	if err != nil {
		log.Errorln("failed to get report for event:", eventId)
		log.Errorln("Error:" , err)
		report.Response.Status = "ERR"
		report.Response.Message = err.Error()
		return report
	}

	report.Response.Status = "OK"
	report.Report = resp
	return report
}
