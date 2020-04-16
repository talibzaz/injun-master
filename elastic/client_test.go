package elastic

import (
	"testing"
	"context"
	"github.com/spf13/viper"
	"fmt"
	"encoding/json"
	"github.com/graphicweave/injun/proto"
	"os"
)

func TestElasticSearch_Search(t *testing.T) {

	viper.Set("ES_HOST", "139.59.85.55")
	viper.Set("ES_PORT", "9200")

	es, _ := NewElasticSearch(context.Background())
	i, err := es.Search("Fre", "#", "#","#",[]string{"be56v7ijbmo0009qnd4g" ,"bdqk4f4mandg00a9ugk0", "bde1rk0dlj0000fct0v0"},false)

	if err != nil {
		t.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(i)


}

func TestElasticSearch_NearbyEvent(t *testing.T) {
	viper.Set("ES_HOST", "139.59.85.55")
	viper.Set("ES_PORT", "9200")

	es, _ := NewElasticSearch(context.Background())
	i, err := es.NearbyEvent( 34.083656, 74.797371)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(i)
}

func TestElasticSearch_UpdateEventById(t *testing.T) {
	viper.Set("ES_HOST", "139.59.85.55")
	viper.Set("ES_PORT", "9200")

	viper.Set("ES_INDEX", "events")
	viper.Set("ES_INDEX_TYPE", "event")



	eventBytes :=
	`{
  "event_id": "bem7jmrohqi0009mpv90",
  "event": {
    "eventDetail": {
			"address": "Baghat, Srinagar",
			"address_2": "barzulla, baghat",
			"allow_brochure_enquires": "yes",
			"allow_exhibitor_enquires": "yes",
			"allow_sponsor_enquires": "no",
			"brief_description": "description",
			"cover_image_thumbnail_upload_id": "shutterstock_302497136_6a674cf1-dd29-4a88-9450-b2bcd2b0b1a9.jpeg",
			"cover_image_upload_id": "shutterstock_302497136_6a674cf1-dd29-4a88-9450-b2bcd2b0b1a9.jpeg",
			"created_on": "1538030040737",
			"detailed_description": "[{\"insert\":\"detailed description\\n\"}]",
			"end_date": "2018-10-02",
			"end_time": "11:28 AM",
			"event_tags": "event",
			"id": "bem7jmrohqi0009mpv90",
			"name": "Test Edit Event",
			"start_date": "2018-09-30",
			"start_time": "2:27 PM",
			"timezone": "GMT+0530",
			"title": "testing edit event",
			"user_id": "6393ffc5-9ba0-429d-aaa3-f938857afdbb",
			"venue_city": "srinagar",
			"venue_country": "India",
			"venue_name": "graphicweaave",
			"venue_state": "JK",
			"zone": "Asia/Calcutta"
		},
		"organizer": {
			"description": "Steel and Wallet.",
			"et_commision_rate": 2,
			"id": "68d50854-0b14-4982-ad38-ebc6063fb0c4",
			"name": "Jambo",
			"status": "APPROVED",
			"upload_id": "2016%2F02%2F22%2F04%2F24%2F31%2Fb7bd820a-ecc0-4170-8f4e-3db2e73b0f4a%2F550250_artsigma_0c3514ae-3b77-408d-8242-c61fc6c70ea4.png",
			"website": "http://www.hakizimana.com"
		},
		"tax": {
			"country_name": "India",
			"should_add_tax": "include",
			"tax_id": "id2",
			"tax_name": "12",
			"tax_rate": "10"
		},
		"ticket": {
			"allow_visitor_registrations": "false",
			"currency": "Canadian Dollar",
			"end_date": "2018-09-30",
			"end_time": "9:52 AM",
			"name": "edited event ticket",
			"price": 20,
			"quantity": 100,
			"sold": 0,
			"start_date": "2018-09-27",
			"start_time": "6:51 PM"
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
		"is_featured": "no"
  }
}`
	var e1 *event.UpdateRequest
	err := json.Unmarshal([] byte(eventBytes), &e1)
	if err != nil {
		t.Logf("could not unmarshal given data", err)
	}

	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", " ")
	//enc.Encode(e1)

	es, _ := NewElasticSearch(context.Background())
	err = es.UpdateEventById(e1)

	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	return
}

func TestElasticSearch_PublishEvent(t *testing.T) {
	viper.Set("ES_HOST", "139.59.85.55")
	viper.Set("ES_PORT", "9200")
	viper.Set("ES_INDEX", "events")
	viper.Set("ES_INDEX_TYPE", "event")

	es, _ := NewElasticSearch(context.Background())
	err := es.PublishEvent( "bf1ifgkru5fg00cvqmcg","12-12-18")

	if err != nil {
		t.Fatal(err)
	}
}