package database

import (
	"testing"
	"github.com/spf13/viper"
	"context"
	"fmt"
	"os"
	"encoding/json"
	"github.com/graphicweave/injun/proto"
	"reflect"
)

func TestArangoDB_RawQuery(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	query := `FOR e in events 
	FILTER e.eventDetail.user_id == @user_id
	RETURN {id: e._key, 
		name: e.eventDetail.name, 
		start_date: e.eventDetail.start_date, 
		end_date: e.eventDetail.end_date, 
		venue_name: e.eventDetail.venue_name, 
		venue_city: e.eventDetail.venue_city, 
		currency: e.ticket.currency, 
		price: e.ticket.price, 
		visitor_registration: 
		e.ticket.allow_visitor_registrations
	}`

	result, err := d.GetEventsRawQuery(query,
		map[string]interface{}{"user_id": "usr-123123"})
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(result[0].Name)
}

func TestArangoDB_GetEventById(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	a, err := NewArangoDB(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	e, err := a.GetEventById("", "bd25vvti78ug00f1o0u0")
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(&e)
}

func TestArangoDB_GetStats(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	a, err := NewArangoDB(context.Background())
	if err != nil {
		t.Fatal(err)
	}


	c, err := a.GetStats("6393ffc5-9ba0-429d-aaa3-f938857afdbb")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
}

func TestArangoDB_GetSold(t *testing.T) {
	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	a, err := NewArangoDB(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	c, err := a.GetSold("bfqhae857r5000d46rd0")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
}

func TestArangoDB_GetWishlist(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	a, err := NewArangoDB(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	w, err := a.GetWishlist("6393ffc5-9ba0-429d-aaa3-f938857afdbb")
	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(" ", " ")
	enc.Encode(&w[0])

}

func TestArangoDB_GetSales(t *testing.T) {

	eventId := "bde1rk0dlj0000fct0v0"

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}
	sales, err := d.GetSales(eventId)

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", " ")
	e.Encode(sales)
}

func TestArangoDB_UpdateEventById(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	eventBytes :=
		`{"event_id": "bem7jmrohqi0009mpv90",
	 "event": {"eventDetail": {
      "id": "bem7jmrohqi0009mpv90",
      "name": "Test Event",
      "brief_description": "description",
      "start_date": "2018-09-27",
      "start_time": "2:27 PM",
      "end_date": "2018-09-28",
      "end_time": "11:28 AM",
      "venue_name": "graphicweaave",
      "address": "Baghat, Srinagar",
      "address_2": "barzulla, baghat",
      "venue_city": "srinagar",
      "venue_country": "India",
      "venue_state": "JK",
      "cover_image_upload_id": "shutterstock_302497136_6a674cf1-dd29-4a88-9450-b2bcd2b0b1a9.jpeg",
      "title": "testing edit event",
      "detailed_description": "[{\"insert\":\"detailed description\\n\"}]",
      "event_tags": "event",
      "user_id": "6393ffc5-9ba0-429d-aaa3-f938857afdbb",
      "timezone": "GMT+0530",
      "zone": "Asia/Calcutta",
      "created_on": "1538030040737",
      "allow_sponsor_enquires": "no",
      "allow_exhibitor_enquires": "yes",
      "allow_brochure_enquires": "yes",
      "cover_image_thumbnail_upload_id": "shutterstock_302497136_6a674cf1-dd29-4a88-9450-b2bcd2b0b1a9.jpeg"
    },
    "organizer": {
      "name": "Jambo",
      "description": "Steel and Wallet.",
      "website": "http://www.hakizimana.com",
      "id": "68d50854-0b14-4982-ad38-ebc6063fb0c4",
      "status": "APPROVED",
      "upload_id": "2016%2F02%2F22%2F04%2F24%2F31%2Fb7bd820a-ecc0-4170-8f4e-3db2e73b0f4a%2F550250_artsigma_0c3514ae-3b77-408d-8242-c61fc6c70ea4.png",
      "et_commision_rate": 2
    },
    "tax": {
      "tax_name": "12",
      "tax_rate": "10",
      "tax_id": "id2",
      "country_name": "India",
      "should_add_tax": "include"
    },
    "ticket": {
      "allow_visitor_registrations": "false",
      "currency": "Canadian Dollar",
      "end_date": "2018-09-28",
      "end_time": "9:52 AM",
      "name": "edited event ticket",
      "price": 20,
      "quantity": 100,
      "start_date": "2018-09-27",
      "start_time": "6:51 PM",
      "sold": 0
    },
    "categories": [
      "Auto & Automotive"
    ],
    "interests": [
      "Communications"
    ],
    "attendees": [
      "Marketing Executives",
      "Director"
    ],
    "eventTypes": [
      "Webinar"
    ],
    "status": "PUBLISHED",
    "coordinates": {
      "lat": 34.0411265,
      "lon": 74.802458
    },
    "mobileApp": {
      "amenities": "NA",
      "help": "NA"
    },
    "is_featured": "yes"
  }}`
	var e1 *event.UpdateRequest

	if err := json.Unmarshal([] byte(eventBytes), &e1); err != nil {
		t.Log("could not unmarshall eventbytes")
	}

	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", " ")
	//enc.Encode(e1)


	db, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	if err = db.UpdateEventById(e1); err != nil {
		t.Log(err)
		t.Fatal(err)
	}

}

func TestArangoDB_UpdateFeaturedEventById(t *testing.T) {

	eventId := "bem7jmrohqi0009mpv90"
	f := "no"

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	err = d.UpdateFeaturedEventById(eventId, f)

	if err != nil {
		t.Fatal(err)
	}

}

func TestUniq(t *testing.T) {
	emails := []string{"abc@mail.com", "abc@mail.com", "abc@mail.com", "xyz@mail.com", "john@doe.com"}

	type response struct {
		Name string
		Address string
		Email string
	}

	resp := response{}

	fmt.Println(reflect.DeepEqual(resp, response{}))

	fmt.Println(uniq(emails))
}

func TestArangoDB_GetFeaturedEvents(t *testing.T) {

	userId := "1288dece-9824-4025-b142-092b3f85dc23"

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	events, err := d.GetFeaturedEvents(userId)

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", " ")
	e.Encode(events)
}

func TestArangoDB_GetGeneralReport(t *testing.T) {

	userid := "1288dece-9824-4025-b142-092b3f85dc23"

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	reports, err := d.GetGeneralReport(userid)

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("",  " ")
	e.Encode(&reports)

}

func TestArangoDB_GetAlternateEvents2(t *testing.T) {
	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")

	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	report, err := d.GetAlternateEvents()

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("",  " ")
	e.Encode(&report)

}

func TestArangoDB_GetReportByEventId(t *testing.T) {

	eventId := "bfctkhpqtoq000ebs0g0"

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")


	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	report, err := d.GetReportByEventId(eventId)

	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("",  " ")
	e.Encode(&report)

}

func TestArangoDB_GetAlternateEvents(t *testing.T) {

	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")
	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	report, err := d.GetAlternateEvents()

	fmt.Println(report)
}

func TestArangoDB_GetWishlisted(t *testing.T) {
	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")
	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	report, err := d.GetWishlisted("1288dece-9824-4025-b142-092b3f85dc23")

	fmt.Println(report)
}


func TestArangoDB_PublishEvent(t *testing.T) {
	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")
	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}
	d.PublishEvent("13524711","12-12-2018")


}

func TestArangoDB_GetEventDetails(t *testing.T) {
	viper.Set("ARANGO_HOST", "http://139.59.85.55:8529")
	viper.Set("ARANGO_DB", "eventackle")
	viper.Set("ARANGO_USERNAME", "root")
	viper.Set("ARANGO_PASSWORD", "qF3mKQcu7zyzBYly")
	d, err := NewArangoDB(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}
	eventDetails, err := d.GetEventDetails("beu5mrfhnklq39f9smsg`", "beu5n0nhnklq39f9smt0")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("",  " ")
	e.Encode(&eventDetails)

}