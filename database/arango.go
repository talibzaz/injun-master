/**
* TODO: Write ArangoDB package DOC
* Assumes Database and all collections exists
**/
package database

import (
	"github.com/arangodb/go-driver/http"
	arango "github.com/arangodb/go-driver"
	"github.com/spf13/viper"
	"context"
	"github.com/graphicweave/injun/proto"
	log "github.com/sirupsen/logrus"
	"strings"
	"github.com/graphicweave/injun/mail"
	"github.com/graphicweave/injun/analytics"
	"fmt"
)

type ArangoDB struct {
	client arango.Client
	conn   arango.Connection
	ctx    context.Context
}

func NewArangoDB(ctx context.Context) (*ArangoDB, error) {

	conn, err := getConnection()
	if err != nil {
		return nil, err
	}

	client, err := arango.NewClient(arango.ClientConfig{
		Connection: conn,
		Authentication: arango.
			BasicAuthentication(viper.GetString("ARANGO_USERNAME"), viper.GetString("ARANGO_PASSWORD")),
	})

	if err != nil {
		return nil, err
	}

	return &ArangoDB{client: client, conn: conn, ctx: ctx}, nil
}

func (a *ArangoDB) Database(db string) (arango.Database, error) {
	if len(db) == 0 {
		db = viper.GetString("ARANGO_DB")
	}
	return a.client.Database(a.ctx, db)
}

func (a *ArangoDB) SaveEvent(event *event.Event) (string, error) {

	db, err := a.Database("")
	if err != nil {
		return "", err
	}

	collection, err := db.Collection(a.ctx, "events")
	if err != nil {
		return "", err
	}

	meta, err := collection.CreateDocument(a.ctx, event)
	if err != nil {
		return "", err
	}

	return meta.Key, nil
}

func (a *ArangoDB) RollbackArango(arangoID string) error {
	db, err := a.Database("")
	if err != nil {
		return err
	}

	collection, err := db.Collection(a.ctx, "events")
	if err != nil {
		return err
	}

	_, err = collection.RemoveDocument(a.ctx, arangoID)
	return err
}

// TODO
func (a *ArangoDB) RollbackUpdateArango(arangoID string) error {
	return nil
}

func (a *ArangoDB) GetEventsRawQuery(query string, bindVars map[string]interface{}) ([]event.EventItem, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	result := make([]event.EventItem, 0)
	for {
		var e event.EventItem

		_, err = cursor.ReadDocument(a.ctx, &e)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return result, nil
}

func (a *ArangoDB) GetEventById(userId, eventId string) (*event.Event, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	bindVars := make(map[string]interface{})

	query := `FOR e in events 
		FILTER e.eventDetail.id == @event_id
		RETURN { 
    		event: e,  
    		is_in_wishlist: (
        		FOR w in wishlists 
        		FILTER w.user_id == @user_id
        		RETURN POSITION(w.events, @event_id) 
    		)[0] ? 'yes': 'no'
		}`

	bindVars["event_id"] = eventId
	bindVars["user_id"] = userId

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	var _event struct {
		Event        event.Event `json:"event"`
		IsInWishlist string      `json:"is_in_wishlist"`
	}

	_, err = cursor.ReadDocument(a.ctx, &_event)

	if err != nil {
		if arango.IsNoMoreDocuments(err) {
			return &event.Event{}, nil
		}
		return nil, err
	}

	_event.Event.IsInWishlist = _event.IsInWishlist

	return &_event.Event, nil
}

func (a *ArangoDB) UpdateEventById(event *event.UpdateRequest) error {

	db, err := a.Database("")
	if err != nil {
		return err
	}

	query := `for e in events
		filter e.eventDetail.id == @id
		replace e with @editedEvent in events`

	bindVars := make(map[string]interface{})
	bindVars["editedEvent"] = event.Event
	bindVars["id"] = event.EventId

	_, err = db.Query(a.ctx, query, bindVars)
	return err
}

func (a *ArangoDB) GetSold(eventId string) (int32, error) {
	db, err := a.Database("")
	if err != nil {
		return 0, err
	}

	query := `for e in events
		filter e.eventDetail.id == @id
		return e.ticket.sold == null ? 0 : e.ticket.sold`


	bindVars := make(map[string]interface{})
	bindVars["id"] = eventId

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var sold int32
	_, err = cursor.ReadDocument(a.ctx, &sold)
	if arango.IsNoMoreDocuments(err) {
		return sold, nil
	} else if err != nil {
		return 0, err
	}

	return sold, nil
}

func (a *ArangoDB) UpdateWishlist(userId, eventId string, shouldAdd bool) error {

	db, err := a.Database("")
	if err != nil {
		return err
	}

	var query string

	bindVars := make(map[string]interface{})
	bindVars["user_id"] = userId
	bindVars["event_id"] = eventId

	if shouldAdd {
		query = `
				UPSERT { user_id: @user_id }
				INSERT { user_id: @user_id, events: TO_ARRAY(@event_id) }
				UPDATE { user_id: @user_id, events: APPEND(OLD.events, @event_id, true) } in wishlists`
	} else {
		query = `
				FOR w in wishlists
    				FILTER w.user_id == @user_id
    				UPDATE w WITH { events : REMOVE_VALUE(w.events, @event_id) } in wishlists
    				RETURN NEW`
	}

	_, err = db.Query(a.ctx, query, bindVars)
	return err
}

func (a *ArangoDB) GetWishlist(userId string) ([]map[string]interface{}, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	query := `
		FOR w IN wishlists
  			FILTER w.user_id == @user_id
  			FOR evtId IN w.events 
    		FOR e IN events 
    			FILTER e.eventDetail.id == evtId 
    			RETURN {
					event: e.eventDetail,
					categories: e.categories,
					ticket: e.ticket,
					status: e.status,
                    deactivated : e.deactivated == null ? false : e.deactivated
				}`

	bindVars := map[string]interface{}{"user_id": userId}

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	documents := make([]map[string]interface{}, 0)
	for {
		var doc map[string]interface{}
		_, err := cursor.ReadDocument(a.ctx, &doc)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Error("wishlist: error while reading document")
			log.Error("Error: " + err.Error())
			continue
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (a *ArangoDB) GetStats(userId string) (map[string]float64, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	bindVars := make(map[string]interface{})

	query := `LET wishes = (
	FOR w in wishlists
	FILTER w.user_id == @user_id
    RETURN LENGTH( w.events))[0]
	
	LET tckts = LENGTH(
    FOR t in tickets
    FILTER t.UserID == @user_id
    RETURN t)

	RETURN { wishes: wishes, tickets: tckts }`

	bindVars["user_id"] = userId

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	stats := make(map[string]float64)
	_, err = cursor.ReadDocument(a.ctx, &stats)

	return stats, err
}

func (a *ArangoDB) ManageEventById(eventId, user_id string) (event.ManageEventResponse, error) {

	db, err := a.Database("")
	if err != nil {
		return event.ManageEventResponse{Status: "ERR"}, err
	}

	bindVars := make(map[string]interface{})
	query := `LET event = (
    			FOR e in events 
    			FILTER e.eventDetail.id == @event_id && e.eventDetail.user_id == @user_id
    			LET evt = e.eventDetail
    			RETURN { 
    				name: evt.name,
    				start_date: evt.start_date,
    				start_time: evt.start_time,
					end_time: evt.end_time,
					end_date: evt.end_date,
    				price: e.ticket.price,
    				quantity: e.ticket.quantity,
    				sold: e.ticket.sold,
    				venue_name: evt.venue_name,
    				venue_city: evt.venue_city,
    				cover_image_upload_id: evt.cover_image_upload_id,
    				ticket_start_date: e.ticket.start_date,
    				ticket_start_time: e.ticket.start_time,
    				ticket_end_time: e.ticket.end_time,
    				ticket_end_date: e.ticket.end_date,
    				time_zone: evt.zone,
					event_status : e.status,
    				currency: e.ticket.currency,
                    deactivated : e.deactivated,
    				allow_visitor_registrations: e.ticket.allow_visitor_registrations == "false" ? 'CLOSED' : 'OPEN'
    			})

			LET visitors = (
			    RETURN SUM(
    			    FOR t in tickets 
    			    FILTER t.EventID == @event_id
    			    RETURN t.NoOfVisitors
    			)
			)[0]

			RETURN MERGE(
     			event[0] ? event[0]: {},
    			{ visitors: visitors }
			)`

	bindVars["event_id"] = eventId
	bindVars["user_id"] = user_id

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return event.ManageEventResponse{Status: "ERR"}, err
	}

	defer cursor.Close()

	var _event event.ManageEventResponse

	_, err = cursor.ReadDocument(a.ctx, &_event)

	if err != nil {
		if arango.IsNoMoreDocuments(err) {
			return event.ManageEventResponse{Status: "ERR"}, err
		}
		return event.ManageEventResponse{Status: "ERR"}, err
	}

	_event.Status = "OK"
	return _event, nil
}

type Sale struct {
	PurchasedOn  int64   `json:"purchased_on"`
	OrderNumber  string  `json:"order_number"`
	TotalAttendees int     `json:"total_attendees"`
	TotalVisitors int 	`json:"total_visitors"`
	Email        string  `json:"email"`
	Position     string  `json:"position"`
	Name         string  `json:"name"`
	AmountPaid   float64 `json:"amount_paid"`
	Attendees []struct {
		Name string `json:"name"`
		Email     string `json:"email"`
	} `json:"attendees"`
}

type Enquiry struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Type           string `json:"type"`
	EventId        string `json:"event_id"`
	JobTitle       string `json:"job_title"`
	Phone          string `json:"phone"`
	Company        string `json:"company"`
	Email          string `json:"email"`
	CompanyWebsite string `json:"company_website"`
}

func (a *ArangoDB) GetSales(eventID string) ([]Sale, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	bindVars := make(map[string]interface{})

	query := `FOR t in tickets
				FILTER t.EventID == @event_id
				RETURN {
    				purchased_on : t.PurchasedOn,
    				order_number : t.ID,
					total_attendees : t.NoOfAttendees,
					total_visitors : t.NoOfVisitors,
    				email: t.Email,
    				attendees: t.Attendees,
					name: t.Name,
					amount_paid: t.AmountCharged
				}`

	bindVars["event_id"] = eventID

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	sales := make([]Sale, 0)

	for {
		var sale Sale

		_, err = cursor.ReadDocument(a.ctx, &sale)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}

	emails := make([]string, 0)
	for k := range sales {
		emails = append(emails, sales[k].Email)
	}

	if len(emails) > 0 {
		uniqEmails := uniq(emails)

		positions, err := GetPositions(len(uniqEmails), uniqEmails)
		if err != nil {
			return nil, err
		}

		for email, position := range positions {
			updateSales(&sales, email, position)
		}
	}

	return sales, nil
}

func uniq(elements []string) []interface{} {

	els := make([]interface{}, 0, len(elements))
	m := make(map[string]bool)

	for _, val := range elements {
		if _, ok := m[val]; !ok {
			m[val] = true
			els = append(els, val)
		}
	}

	return els
}

func updateSales(sales *[]Sale, email, position string) {
	for k := range *sales {
		if (*sales)[k].Email == email {
			(*sales)[k].Position = position
		}
	}
}

func (a *ArangoDB) GetWishlisted(userId string) ([]string, error) {
	db, err := a.Database("")
	if err != nil {
		return nil, err
	}
	bindVars := make(map[string]interface{})

	query := `for w in wishlists
					filter w.user_id == @user_id
						return w.events`
	bindVars["user_id"] = userId
	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	var events []string
	_, err = cursor.ReadDocument(a.ctx, &events)
	if arango.IsNoMoreDocuments(err) {
		return events, nil
	} else if err != nil {
		return nil, err
	}

	return events, nil

}

func (a *ArangoDB) UpdateFeaturedEventById(eventId, featured string) error {

	db, err := a.Database("")
	if err != nil {
		return err
	}

	bindVars := make(map[string]interface{})

	query := `FOR e in events 
				FILTER e.eventDetail.id == @event_id && e.status == "PUBLISHED"
				UPDATE e WITH {is_featured : @is_featured} in events
				RETURN NEW`

	bindVars["event_id"] = eventId
	bindVars["is_featured"] = featured

	_, err = db.Query(a.ctx, query, bindVars)
	if err != nil {
		return err
	}
	return nil

}

func (a *ArangoDB) GetFeaturedEvents(userId string) ([]map[string]interface{}, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	bindVars := map[string]interface{}{"user_id": strings.TrimSpace(userId)}

	query := `FOR e in events 
				FILTER e.is_featured == 'yes'  && DATE_NOW() < DATE_TIMESTAMP(e.eventDetail.end_date) && e.deactivated != true
				LIMIT 10
				RETURN {
    				id: e.eventDetail.id,
    				name: e.eventDetail.name,
    				start_date: CONCAT(e.eventDetail.start_date, ' ', e.eventDetail.start_time),
    				categories: CONCAT_SEPARATOR(', ', e.categories),
					image_id: e.eventDetail.cover_image_thumbnail_upload_id,
    				is_wishlisted: (
        				!@user_id ? false : (
            				FOR w in wishlists
            					FILTER w.user_id == @user_id
            					RETURN w.events[*]
        					)[0] ANY == e.eventDetail.id
    				)
				}`

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	featuredEvents := make([]map[string]interface{}, 0)

	for {
		var evt map[string]interface{}

		_, err = cursor.ReadDocument(a.ctx, &evt)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		featuredEvents = append(featuredEvents, evt)
	}

	return featuredEvents, nil
}

func (a *ArangoDB) GetGeneralReport(userId string) ([]map[string]interface{}, error) {

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	bindVars := map[string]interface{}{"user_id": strings.TrimSpace(userId)}

	query := `LET reports = (
    	FOR e in events 
        	FILTER e.eventDetail.user_id == @user_id && e.status == "PUBLISHED"
        	COLLECT ticket = e.ticket.sold, event = e.eventDetail INTO report_groups = {
            	"id": e.eventDetail.id,
            	"name" :  e.eventDetail.name,
            	"start_date" :  e.eventDetail.start_date,
				"end_date" : e.eventDetail.end_date,
            	"venue": CONCAT(e.eventDetail.venue_name, ', ', e.eventDetail.venue_city),
            	"tickets": e.ticket.quantity,
            	"sold" : e.ticket.sold == null ? 0 :e.ticket.sold,
            	"visitors": (
            	        return sum (
            	        for t in tickets
            	        filter t.EventID ==  e.eventDetail.id
            	        return t.NoOfVisitors
            	    ))[0],
            	"revenue": e.ticket.sold == 0 ? 0 : (
					RETURN SUM (
						FOR t in tickets
						FILTER t.EventID == e.eventDetail.id
						RETURN t.AmountCharged/t.ExchangeRate
					)
				)[0]
        	}
        	RETURN {
            	sold: ticket,
            	events: report_groups
        	}
		)
		RETURN {
		    sold: SUM(reports[*].sold),
   			events: reports[*].events
		}`

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	reports := make([]map[string]interface{}, 0)

	for {
		var evt map[string]interface{}

		_, err = cursor.ReadDocument(a.ctx, &evt)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		reports = append(reports, evt)
	}

	return reports, nil

}

func (a *ArangoDB) GetReportByEventId(eventId string) (map[string]interface{}, error) {
	pageViews, err := analytics.GetPageViewsByEventID(eventId)
	if err != nil {
		pageViews = 0
	}

	db, err := a.Database("")
	if err != nil {
		return nil, err
	}

	bindVars := map[string]interface{}{"event_id": strings.TrimSpace(eventId),"page_views": pageViews}

	query := `LET report = (
    	FOR e in events 
        	FILTER e.eventDetail.id == @event_id
        	COLLECT ticket = e.ticket.sold == null ? 0 : e.ticket.sold, event = e.eventDetail INTO report_group = {
            	"id": e.eventDetail.id,
            	"name" :  e.eventDetail.name,
            	"start_date" :  e.eventDetail.start_date,
            	"end_date" :  e.eventDetail.end_date,
				"ticket_start_date": e.ticket.start_date,
				"ticket_end_date" : e.ticket.end_date,
            	"tickets": e.ticket.quantity,
            	"visitors": (
            	        return sum (
            	        for t in tickets
            	        filter t.EventID ==  e.eventDetail.id
            	        return t.NoOfVisitors
            	    ))[0],
				"brochure_requests" : (
			            RETURN LENGTH(
			                FOR b in brochure_requests
			                FILTER b.eventId == @event_id
			                RETURN b
			            )
			                    
			        )[0],
			    "exhibitor_enquiries" : (
			            RETURN LENGTH(
			                FOR b in enquiries
			                FILTER b.eventId == @event_id && b.enquiryType == "exhibitor"
			                RETURN b
			            )
			                    
			        )[0],
			    "sponsor_enquiries": (
			            RETURN LENGTH(
			                FOR b in enquiries
			                FILTER b.eventId == @event_id && b.enquiryType == "sponsor"
			                RETURN b
			             )
			        )[0],
            	"revenue": e.ticket.sold == 0 ? 0 : (
					RETURN SUM (
						FOR t in tickets
						FILTER t.EventID == @event_id
						RETURN t.AmountCharged/t.ExchangeRate
					)
				)[0],
				page_views: @page_views
        	}
        RETURN MERGE(
            report_group[0] ? report_group[0] : {} , {"sold": ticket}
        )
	)
	LET ticket_sale = (
    	FOR t in tickets 
    	FILTER t.EventID == @event_id
		SORT t.PurchasedOn ASC
    	RETURN {purchased_on: DATE_FORMAT(t.PurchasedOn, '%dd-%mm-%yyyy'), sale: t.NoOfAttendees}
	)

	RETURN {report: report[0] ? report[0] : {}, sale: ticket_sale}`

	cursor, err := db.Query(a.ctx, query, bindVars)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	report := make(map[string]interface{})
	_, err = cursor.ReadDocument(a.ctx, &report)

	return report, err

}

func (a *ArangoDB) GetAlternateEvents() ([]map[string]interface{}, error) {
	db, err := a.Database("")
	if err != nil {
		return nil, err
	}
	query := `FOR e IN events 
				filter DATE_NOW() < DATE_TIMESTAMP(e.eventDetail.end_date) && e.status == "PUBLISHED" && e.deactivated != true
				SORT e.eventDetail.created_on DESC
               LIMIT 10
             return {
                "id" :   e.eventDetail.id,
				"name" : e.eventDetail.name,
				"image": e.eventDetail.cover_image_thumbnail_upload_id,
				"start_date":e.eventDetail.start_date,
				"start_time" : e.eventDetail.start_time,
				"categories":  e.categories
                     }`
	cursor, err := db.Query(a.ctx, query, nil)
	if err != nil {
		fmt.Println("RR",err)
		return nil, err
	}

	defer cursor.Close()
	events := make([]map[string]interface{}, 0)

	for {
		var evt map[string]interface{}

		_, err = cursor.ReadDocument(a.ctx, &evt)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		events = append(events, evt)
	}

	return events, nil
}

func (a *ArangoDB) InsertEnquiryDetails(enquiry *mail.EnquiryRequest) error {
	db, err := a.Database("")
	if err != nil {
		return err
	}

	col, err := db.Collection(a.ctx, "enquiries")
	if err != nil {
		return err
	}

	_, err = col.CreateDocument(a.ctx, enquiry)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArangoDB) BrochureRequest(brochure *mail.BrochureRequest) error {
	db, err := a.Database("")
	if err != nil {
		return err
	}

	col, err := db.Collection(a.ctx, "brochure_requests")
	if err != nil {
		return err
	}

	_, err = col.CreateDocument(a.ctx, brochure)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArangoDB) 	PublishEvent(eventID, date string) error {
	db, err := a.Database("")

	col, err := db.Collection(a.ctx, "events")
	if err != nil {
       return err
	}

	patch := map[string]interface{}{
		"status":     "PUBLISHED",
		"eventDetail": map[string]interface{} {"created_on": date},
	}
	_, err = col.UpdateDocument(a.ctx, eventID, patch)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArangoDB) GetEventDetails(ticketId, userId string) (map[string]interface{}, error) {
	db, err := a.Database("")
	if err != nil {
		return nil, err
	}
	var eventDetails map[string]interface{}
	if userId == "" {
		query := `FOR t IN tickets
	FILTER t.ID == @ticketId
    let event = (
    FOR e IN events
    FILTER t.EventID == e.eventDetail.id 
    return{
        eventId : e.eventDetail.id,
        eventName: e.eventDetail.name,
        eventCoverImage: e.eventDetail.cover_image_upload_id,
        eventDateTime: CONCAT(e.eventDetail.start_date," ",e.eventDetail.start_time),
        eventFullVenue: CONCAT(e.eventDetail.venue_name, ", ", e.eventDetail.venue_city, ", ", e.eventDetail.venue_country),
        eventOrganizerCompany: e.organizer.name,
		coordinates: e.coordinates
    })
    return {
        event: event
    }`
		bindVars := map[string]interface{}{"ticketId": ticketId}
		cursor, err := db.Query(a.ctx, query, bindVars)
		if err != nil {
			return nil, err
		}

		defer cursor.Close()

		_, err = cursor.ReadDocument(a.ctx, &eventDetails)

		if err != nil {
			if arango.IsNoMoreDocuments(err) {
				return eventDetails, nil
			}
			return nil, err
		}

	} else {
		query := `FOR t IN tickets
	FILTER t.ID == @ticketId
    let user = (
        FOR a in t.Attendees
        FILTER a.id == @userId
        return {
            name: a.name,
            email: a.email
    })
    let event = (
    FOR e IN events
    FILTER t.EventID == e.eventDetail.id 
    return{
        eventId : e.eventDetail.id,
        eventName: e.eventDetail.name,
        eventCoverImage: e.eventDetail.cover_image_upload_id,
        eventDateTime: CONCAT(e.eventDetail.start_date," ",e.eventDetail.start_time),
        eventFullVenue: CONCAT(e.eventDetail.venue_name, ", ", e.eventDetail.venue_city, ", ", e.eventDetail.venue_country),
        eventOrganizerCompany: e.organizer.name,
		coordinates: e.coordinates
    })
    return {
        event: event,
        user: user
    }`
		bindVars := map[string]interface{}{"ticketId": ticketId, "userId": userId}
		cursor, err := db.Query(a.ctx, query, bindVars)
		if err != nil {
			return nil, err
		}

		defer cursor.Close()

		_, err = cursor.ReadDocument(a.ctx, &eventDetails)

		if err != nil {
			if arango.IsNoMoreDocuments(err) {
				return eventDetails, nil
			}
			return nil, err
		}
	}

	return eventDetails, nil
}

func getConnection() (arango.Connection, error) {
	return http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{viper.GetString("ARANGO_HOST")},
	})
}
