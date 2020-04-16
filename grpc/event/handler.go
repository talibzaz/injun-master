package event

import (
	"golang.org/x/net/context"
	"github.com/graphicweave/injun/proto"
	log "github.com/sirupsen/logrus"
	"github.com/graphicweave/injun/service"
	"errors"
	"github.com/rs/xid"
	"fmt"
	"github.com/graphicweave/injun/analytics"
)

type EventService struct {
}

func (es EventService) UpdateFeaturedEventById(ctx context.Context, evt *event.UpdateFeaturedRequest) (*event.UpdateResponse, error) {
	arangoDb, err := service.NewArangoDB(context.Background())

	log.Infoln("recvd request for event id: ", evt.EventId)

	if err != nil {
		log.Error("failed to get ArangoDB connection")
		log.Error("Error: ", err)
		return &event.UpdateResponse{Status: "ERR"}, err
	}

	err = arangoDb.UpdateFeaturedEventById(evt.EventId, evt.Featured)

	if err != nil {
		log.Error("failed to update featured event by event ID")
		log.Error("event ID: ", evt.EventId)
		log.Error("Featured: ", evt.Featured)
		log.Error("Error: ", err)
		return &event.UpdateResponse{Status: "ERR"}, err
	}

	return &event.UpdateResponse{Status: "OK"}, nil
}

func (es EventService) CreateEvent(ctx context.Context, evt *event.Event) (*event.Response, error) {

	arangoDb, err := service.NewArangoDB(context.Background())
	fmt.Println(evt)
	if err != nil {
		log.Errorln("failed to create ArangoDB client")
		log.Errorln("Error: ", err.Error())
		return &event.Response{ArangoID: "", ElasticId: "", Status: "ERR"}, err
	}

	elasticDb, err := service.NewElasticSearch(ctx)

	if err != nil {
		log.Errorln("failed to create ElasticSearch client")
		log.Errorln("Error: ", err.Error())
		return &event.Response{ArangoID: "", ElasticId: "", Status: "ERR"}, err
	}

	var arangoId, elasticId string
	arangoErrChan := make(chan error)
	elasticErrChan := make(chan error)

	evt.EventDetail.Id = xid.New().String()
	evt.Ticket.Sold = 0

	fmt.Println("as",evt)

	go func() {
		arangoId, err = arangoDb.SaveEvent(evt)
		arangoErrChan <- err
	}()

	go func() {
		elasticId, err = elasticDb.SaveEvent(evt)
		elasticErrChan <- err
	}()

	arangoErr := <-arangoErrChan
	esError := <-elasticErrChan

	if arangoErr != nil && esError == nil {
		log.Error("failed to save in ArangoDB")
		log.Error("Error ", arangoErr.Error())
		go func() {
			err := elasticDb.RollbackES(arangoId)
			if err != nil {
				log.Error("failed to rollback ArangoDB for id", arangoId)
				log.Errorln("Error: ", err.Error())
				return
			}
		}()
		return &event.Response{ArangoID: "", ElasticId: "", Status: "ERR"}, nil
	}

	if arangoErr == nil && esError != nil {
		log.Error("failed to save in ElasticSearch")
		log.Error("Error ", esError)
		go func() {
			err := arangoDb.RollbackArango(arangoId)
			if err != nil {
				log.Error("failed to rollback ArangoDB for id", arangoId)
				log.Errorln("Error: ", err)
				return
			}
		}()
		return &event.Response{ArangoID: "", ElasticId: "", Status: "ERR"}, nil
	}

	if arangoErr != nil && esError != nil {

		log.Error("failed to save event in ArangoDB and ElasticSearch")
		log.Error("arango error ", arangoErr)
		log.Error("elastic search error ", esError)

		return &event.Response{ArangoID: "", ElasticId: "", Status: "ERR"}, nil
	}

	log.Info("event created with Arango ID: ", arangoId)
	log.Info("event created with ElasticSearch with ID: ", elasticId)
	return &event.Response{ArangoID: arangoId, ElasticId: elasticId, Status: "OK"}, nil
}

func (es EventService) Rollback(ctx context.Context, request *event.RollbackRequest) (*event.RollbackReponse, error) {

	arangoId := request.ArangoId
	elasticId := request.ElasticId

	arangoErrChan := make(chan error)
	elasticErrChan := make(chan error)

	arangoDb, err := service.NewArangoDB(context.Background())

	if err != nil {
		log.Errorln("failed to create ArangoDB client")
		log.Errorln("Error: ", err.Error())
		return &event.RollbackReponse{Status: "ERR"}, err
	}

	elasticDb, err := service.NewElasticSearch(ctx)

	if err != nil {
		log.Errorln("failed to create ElasticSearch client")
		log.Errorln("Error: ", err.Error())
		return &event.RollbackReponse{Status: "ERR"}, err
	}

	go func() {
		elasticErrChan <- elasticDb.RollbackES(elasticId)
	}()

	go func() {
		arangoErrChan <- arangoDb.RollbackArango(arangoId)
	}()

	arangoErr := <-arangoErrChan
	esError := <-elasticErrChan

	if arangoErr != nil && esError == nil {
		log.Error("failed to rollback in ArangoDB")
		log.Error("Error ", arangoErr)

		return &event.RollbackReponse{Status: "ERR"}, arangoErr
	}

	if arangoErr == nil && esError != nil {
		log.Error("failed to rollback in ElasticSearch")
		log.Error("Error ", esError)

		return &event.RollbackReponse{Status: "ERR"}, esError
	}

	if arangoErr != nil && esError != nil {

		log.Error("failed to rollback event in ArangoDB and ElasticSearch")
		log.Error("arango error ", arangoErr)
		log.Error("elastic search error ", esError)

		return &event.RollbackReponse{Status: "ERR"}, errors.New(arangoErr.Error() + "||" + esError.Error())
	}

	log.Info("event rollbacked with Arango ID: ", arangoId)
	log.Info("event rollbacked with ElasticSearch with ID: ", elasticId)
	return &event.RollbackReponse{Status: "OK"}, nil
}

func (es EventService) GetEventsByUserId(evt *event.EventRequest, stream event.EventService_GetEventsByUserIdServer) error {

	arangoDb, err := service.NewArangoDB(context.Background())
	if err != nil {
		log.Error("failed to get ArangoDB connection")
		log.Error("Error: ", err)
		return err
	}

	query := `FOR e in events 
	FILTER e.eventDetail.user_id == @user_id
	RETURN {id: e.eventDetail.id, 
		name: e.eventDetail.name, 
		start_date: e.eventDetail.start_date, 
		end_date: e.eventDetail.end_date, 
		venue_name: e.eventDetail.venue_name, 
		venue_city: e.eventDetail.venue_city, 
		currency: e.ticket.currency, 
		price: e.ticket.price, 
		visitor_registration: 
		e.ticket.allow_visitor_registrations,
		cover_image_upload_id: e.eventDetail.cover_image_thumbnail_upload_id,
        status: e.status,
        start_time: e.eventDetail.start_time,
        end_time: e.eventDetail.end_time,
        zone:e.eventDetail.zone,
        timezone: e.eventDetail.timezone,
		tickets: e.ticket.quantity,
		sold: e.ticket.sold,
        key: e._key,
        organizer_id : e.organizer.id,
        organizer_name :  e.organizer.name,
		organizer_status : e.organizer.status,
		deactivated : e.deactivated
	}`
	bindVars := map[string]interface{}{
		"user_id": evt.UserId,
	}

	events, err := arangoDb.GetEventsRawQuery(query, bindVars)

	if err != nil {
		log.Error("failed to query arango for user events")
		log.Error("user id: ", evt.UserId)
		log.Error("Error: ", err)
		return err
	}
	for k := range events {
		if err := stream.Send(&events[k]); err != nil {
			log.Error("failed to send event to grpc stream")
			log.Error("Error: ", err)
			return err
		}
	}
	log.Info("event items streamed!")

	return nil
}

func (es EventService) GetEventById(ctx context.Context, evt *event.EventRequest) (*event.Event, error) {

	arangoDb, err := service.NewArangoDB(context.Background())

	log.Infoln("recvd request for event id: ", evt.EventId)
	log.Infoln("recvd request for user id: ", evt.UserId)

	if err != nil {
		log.Error("failed to get ArangoDB connection")
		log.Error("Error: ", err)
		return nil, err
	}

	_event, err := arangoDb.GetEventById(evt.UserId, evt.EventId)

	if err != nil {

		log.Error("failed to get event by event ID")
		log.Error("event ID: ", evt.EventId)
		log.Error("Error: ", err)
		return nil, err
	}

	return _event, nil
}

func (es EventService) UpdateEventById(ctx context.Context, evt *event.UpdateRequest) (*event.UpdateResponse, error) {

	arangoDb, err := service.NewArangoDB(context.Background())
	if err != nil {
		log.Error("failed to get ArangoDB connection")
		log.Error("Error: ", err)
		return &event.UpdateResponse{Status: "ERR"}, err
	}

	elasticDb, err := service.NewElasticSearch(ctx)

	if err != nil {
		log.Errorln("failed to get ElasticSearch client")
		log.Errorln("Error: ", err.Error())
		return &event.UpdateResponse{ Status: "ERR"}, err
	}

	var arangoId, elasticId string
	arangoErrChan := make(chan error)
	elasticErrChan := make(chan error)

	go func() {
		evt.Event.Ticket.Sold, err =  arangoDb.GetSold(evt.EventId)
		arangoErrChan <- arangoDb.UpdateEventById(evt)
	}()

	go func() {
		elasticErrChan <- elasticDb.UpdateEventById(evt)
	}()

	arangoErr := <-arangoErrChan
	esError := <-elasticErrChan

	if arangoErr != nil && esError == nil {
		log.Error("failed to update in ArangoDB")
		log.Error("Error ", arangoErr.Error())
		go func() {
			err := elasticDb.RollbackUpdateES(arangoId)
			if err != nil {
				log.Error("failed to rollback ArangoDB for id", arangoId)
				log.Errorln("Error: ", err)
				return
			}
		}()
		return &event.UpdateResponse{Status: "ERR"}, nil
	}

	if arangoErr == nil && esError != nil {
		log.Error("failed to update in ElasticSearch")
		log.Error("Error ", esError)
		go func() {
			err := arangoDb.RollbackUpdateArango(arangoId)
			if err != nil {
				log.Error("failed to rollback ArangoDB for id", arangoId)
				log.Errorln("Error: ", err)
				return
			}
		}()
		return &event.UpdateResponse{Status: "ERR"}, nil
	}

	if arangoErr != nil && esError != nil {

		log.Error("failed to update event in ArangoDB and ElasticSearch")
		log.Error("arango error ", arangoErr)
		log.Error("elastic search error ", esError)

		return &event.UpdateResponse{Status: "ERR"}, nil
	}

	log.Infof("event with Arango ID %s updated ", arangoId)
	log.Infof("event with ElasticID %s updated ", elasticId)
	return &event.UpdateResponse{Status: "OK"}, nil
}

func (es EventService) ManageEventById(ctx context.Context, evt *event.EventRequest) (*event.ManageEventResponse, error) {
	arangoDb, err := service.NewArangoDB(context.Background())
	if err != nil {
		log.Error("failed to get ArangoDB connection")
		log.Error("Error: ", err)
		return &event.ManageEventResponse{Status: "ERR"}, err
	}

	_event, err := arangoDb.ManageEventById(evt.EventId, evt.UserId)
	if err != nil {
		log.Error("failed to update document for event")
		log.Error("event_id: ", evt.EventId)
		log.Error("Error: ", err)
		return &event.ManageEventResponse{Status: "ERR"}, err
	}
	_event.PageViews, err = analytics.GetPageViewsByEventID(evt.EventId)
	return &_event, nil
}
