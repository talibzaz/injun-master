package elastic

import (
	"github.com/olivere/elastic"
	"context"
	"github.com/spf13/viper"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/graphicweave/injun/proto"
	"encoding/json"
	"text/template"
	"sync"
	"bytes"
	"errors"
	"strings"
	"time"
	currency "github.com/younisshah/go-currency-code"
)

var ErrEsNoResults = errors.New("no results found")

type ElasticSearch struct {
	addr   string
	client *elastic.Client
	ctx    context.Context
}

var (
	localEventsQueryTemplate  *template.Template
	locationOnlyQueryTemplate *template.Template
	countryOnlyQueryTemplate  *template.Template
	stateOnlyQueryTemplate    *template.Template
	cityOnlyQueryTemplate     *template.Template
	termAndCityQueryTemplate  *template.Template
	termAndStateQueryTemplate  *template.Template
	termAndCountryQueryTemplate  *template.Template
	searchQueryTemplate       *template.Template
	once                      sync.Once
)

const mapping = `
	{
   "mappings": {
       "event": {
           "properties": {
				"coordinates": {
                    "type": "geo_point"
				}
			}
       }
   }
}
`

func NewElasticSearch(context context.Context) (*ElasticSearch, error) {

	client, err := getClient()

	if err != nil {
		return nil, err
	}

	es := &ElasticSearch{
		addr:   "http://" + viper.GetString("ES_HOST") + ":" + viper.GetString("ES_PORT"),
		client: client,
		ctx:    context,
	}

	// compile the query once
	once.Do(func() {

		localEventsQueryTemplate, _ = template.New("es_local_query_template").
			Parse(`{ 
                             "bool":{ "must":{ "query_string":{ "query":"*" } } ,
      							  "filter": { "geo_distance": { "distance": "{{.Distance}}km",
                                                                "coordinates": {
                                                                                 "lat": {{.Lat}},
                                                                                 "lon": {{.Lon}}
                                                                                }
     												          }
                                            }
       							}
                       }`)

		searchQueryTemplate, _ = template.New("es_query_template").
			Parse(`{"query_string": {"query": "{{.}}*"}}`)

		countryOnlyQueryTemplate, _ = template.New("country_query_template").
			Parse(`{
                        "bool": {
                               "must": {
                                   "match": {
                                         "eventDetail.venue_country": "{{.Country}}"
                                           }
                                       }
                                }
                    }`)
		stateOnlyQueryTemplate, _ = template.New("state_query_template").
			Parse(`{
                        "bool": {
                               "must": {
                                   "match": {
                                         "eventDetail.venue_state": "{{.State}}"
                                           }
                                       }
                                }
                    }`)
		cityOnlyQueryTemplate, _ = template.New("state_query_template").
			Parse(`{
                        "bool": {
                               "must": {
                                   "match": {
                                          "eventDetail.venue_city": "{{.City}}"
                                           }
                                       }
                                }
                    }`)

		locationOnlyQueryTemplate, _ = template.New("es_location_only_query_template").
			Parse(`
							{
    "bool": {
      "must": [
        {
          "match": {
            "eventDetail.venue_city": "{{.City}}"
          }
        },
        {
          "match": {
            "eventDetail.venue_state": "{{.State}}"
          }
        },
        {
          "match": {
            "eventDetail.venue_country": "{{.Country}}"
          }
        }
      ]
    }
  }`)

		termAndCityQueryTemplate, _ = template.New("es_term_city_only_query_template").
			Parse(`{

     "bool": {
       "must": [
         {
             "query_string": {
             "query": "{{.Term}}*"
          }
        },
        {
          "match": {
            "eventDetail.venue_city": "{{.City}}"
          }
        }
      ]
    
  }
}`)
		termAndCountryQueryTemplate, _ = template.New("es_term_country_only_query_template").
			Parse(`{
 
     "bool": {
       "must": [
         {
             "query_string": {
             "query": "{{.Term}}*"
          }
        },
        {
          "match": {
            "eventDetail.venue_country": "{{.Country}}"
          }
        }
      ]
    
  }
}`)
		termAndStateQueryTemplate, _ = template.New("es_term_state_only_query_template").
			Parse(`{

     "bool": {
       "must": [
         {
             "query_string": {
             "query": "{{.Term}}*"
          }
        },
        {
          "match": {
            "eventDetail.venue_state": "{{.State}}"
          }
        }
      ]
    }
  
}`)

	})

	return es, nil
}

func (es *ElasticSearch) Ping() error {

	_, _, err := es.client.Ping(es.addr).Do(es.ctx)
	return err
}

func (es *ElasticSearch) SetupIndex() error {

	index := viper.GetString("ES_INDEX")

	exists, err := es.client.IndexExists(index).Do(es.ctx)
	if err != nil {
		return err
	}
	if !exists {
		createIndex, err := es.client.CreateIndex(index).BodyString(mapping).Do(es.ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return fmt.Errorf("failed to get acknowledgement for index creation " + index)
		}
		log.Infoln("index " + index + " successfully created")
		return nil

	} else {
		log.Infoln("index " + index + " already exists")
		return nil
	}
}

type SearchResponse struct {
	SearchResults []SearchResult
	SearchItems   []SearchItem
}

// SearchResult is one ElasticSearch search result
type SearchResult struct {
	Event struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"eventDetail"`
}

type SearchItem struct {
	EventDetail struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Venue        string `json:"venue_name"`
		VenueCity    string `json:"venue_city"`
		VenueCountry string `json:"venue_country"`
		StartDate    string `json:"start_date"`
		StartTime    string `json:"start_time"`
		EndTime      string `json:"end_time"`
		EndDate		 string	`json:"end_date"`
		Zone	 	 string	`json:"zone"`
		CoverImage   string `json:"cover_image_thumbnail_upload_id"`
	} `json:"eventDetail"`
	Ticket struct {
		Currency string  `json:"currency"`
		Price    float64 `json:"price"`
	} `json:"ticket"`
	Categories []string `json:"categories"`
	Interests  []string `json:"interests"`
	Attendees  []string `json:"attendees"`
	EventTypes []string `json:"eventTypes"`
	IsWishlisted bool   `json:"is_wishlisted"`
	Status 		string  `json:"status"`
	Deactivated bool 	`json:"deactivated"`
}

func (es *ElasticSearch) Search(searchTerm, country, state, city string,wishlistedEvents []string, isAuto bool) (interface{}, error) {

	//TODO refactor location Search.
	var searchResults *elastic.SearchResult
	var err error

	if searchTerm != "#" {
		if country != "#" || state != "#" || city != "#" {
			searchResults, err = es.terms(strings.ToLower(country), strings.ToLower(state), strings.ToLower(city) ,strings.ToLower(searchTerm))
		} else {
			searchResults, err = es.search(strings.ToLower(searchTerm))
		}
	} else {
		searchResults, err = es.locationOnly(strings.ToLower(country), strings.ToLower(state), strings.ToLower(city))
	}

	if err != nil {
		return SearchResponse{}, err
	}
	log.Infof("got %d hits", searchResults.Hits.TotalHits)

	if searchResults.Hits.TotalHits > 0 {

		response := SearchResponse{
			SearchItems:   make([]SearchItem, 0),
			SearchResults: make([]SearchResult, 0),
		}

		for _, v := range searchResults.Hits.Hits {

			resultBytes, err := v.Source.MarshalJSON()
			if err != nil {
				log.Error("failed to marshal JSON search result")
				log.Errorln("Error ", err)
				return nil, err
			}


			var event SearchItem
			json.Unmarshal(resultBytes, &event)
			if err != nil {
				log.Error("failed to unmarshal JSON search result to SearchResult")
				log.Errorln("Error ", err)
				continue
			}
			if event.Status == "DRAFT" {
				fmt.Println("skipping draft event.",event.EventDetail.Name)
				continue
			}

			if event.Deactivated == true {
				log.Info("skipping deactivated event", event.EventDetail.Name)
				continue
			}

			if isPastEvent(event) {
				fmt.Println("skipping past event.", event.EventDetail.Name)
				continue
			}


			switch isAuto {
			case true:
				var r SearchResult

				err = json.Unmarshal(resultBytes, &r)
				if err != nil {
					log.Error("failed to unmarshal JSON search result to SearchResult")
					log.Errorln("Error ", err)
					continue
				}
				fmt.Println(r)
				response.SearchResults = append(response.SearchResults, r)
			case false:
				var r SearchItem
				err = json.Unmarshal(resultBytes, &r)

				for _, v := range wishlistedEvents {
					if r.EventDetail.Id == v {
						r.IsWishlisted = true
					}
				}

				if err != nil {
					log.Error("failed to unmarshal JSON search result to SearchItem")
					log.Errorln("Error ", err)
					continue
				}
				sym , _ := currency.FromCurrencyName(r.Ticket.Currency)

				r.Ticket.Currency = sym["symbol"].(string)
				response.SearchItems = append(response.SearchItems, r)
			}
		}

		return response, err
	}
	return SearchResponse{}, ErrEsNoResults
}

func (es *ElasticSearch) NearbyEvent(lat, lon float64) (interface{}, error) {
	searchResults, err := es.localEvent(lat, lon)

	if err != nil {

		return SearchResponse{}, err
	}
	if searchResults.Hits.TotalHits > 0 {

		searchItems := make([]SearchItem, 0)
		events := make([]map[string]interface{}, 0)
		for _, v := range searchResults.Hits.Hits {

			resultBytes, err := v.Source.MarshalJSON()

			if err != nil {
				log.Error("failed to marshal JSON search result")
				log.Errorln("Error ", err)
				return nil, err
			}

			var r SearchItem

			err = json.Unmarshal(resultBytes, &r)
			if err != nil {
				log.Error("failed to unmarshal JSON search result to SearchItem")
				log.Errorln("Error ", err)
				continue
			}
			if r.Status == "DRAFT" {
				fmt.Println("skipping draft event.",r.EventDetail.Name)
				continue
			}

			if r.Deactivated == true {
				log.Info("skipping deactivated event", r.EventDetail.Name)
				continue
			}


			if isPastEvent(r) {
				fmt.Println("skipping past event.",r.EventDetail.Name)
				continue
			}
			searchItems = append(searchItems, r)
			events = append(events, map[string]interface{}{
				"id":         r.EventDetail.Id,
				"name":       r.EventDetail.Name,
				"image":      r.EventDetail.CoverImage,
				"start_date": r.EventDetail.StartDate,
				"start_time": r.EventDetail.StartTime,
				"categories": r.Categories,
			})
		}
		return events, nil
	}

	return SearchResponse{}, ErrEsNoResults

}

func (es *ElasticSearch) search(searchTerm string) (*elastic.SearchResult, error) {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	searchQueryTemplate.Execute(writer, searchTerm)

	query := elastic.NewRawStringQuery(writer.String())
	return es.client.Search().
		Index(viper.GetString("ES_INDEX")).
		From(0).
		Size(1000).
		Query(query).
		Pretty(true).
		Do(es.ctx)
}

func (es *ElasticSearch) locationOnly(country, state, city string) (*elastic.SearchResult, error) {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)
	data := struct {
		State   string
		Country string
		City    string
	}{
		State:   state,
		Country: country,
		City:    city,
	}

	if city != "#" {
		if state != "#" && country != "#" {
			locationOnlyQueryTemplate.Execute(writer, data)
		} else {
			cityOnlyQueryTemplate.Execute(writer, data)
		}
	} else if state != "#" {
		stateOnlyQueryTemplate.Execute(writer, data)
	} else {
		countryOnlyQueryTemplate.Execute(writer, data)
	}

	query := elastic.NewRawStringQuery(writer.String())

	return es.client.Search().
		Index(viper.GetString("ES_INDEX")).
		From(0).
		Size(10).
		Query(query).
		Pretty(true).
		Do(es.ctx)
}

func (es *ElasticSearch) terms(country, state, city, term string) (*elastic.SearchResult, error) {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	data := struct {
		State   string
		Country string
		City    string
		Term    string
	}{
		State:   state,
		Country: country,
		City:    city,
		Term:    term,
	}

	if city != "#" {
		termAndCityQueryTemplate.Execute(writer, data)
	} else if state != "#" {
		termAndStateQueryTemplate.Execute(writer, data)
	} else {
		termAndCountryQueryTemplate.Execute(writer, data)
	}
	query := elastic.NewRawStringQuery(writer.String())

	return es.client.Search().
		Index(viper.GetString("ES_INDEX")).
		From(0).
		Size(10).
		Query(query).
		Pretty(true).
		Do(es.ctx)
}

func (es *ElasticSearch) localEvent(lat, lon float64) (*elastic.SearchResult, error) {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)
	data := struct {
		Distance int32
		Lat      float64
		Lon      float64
	}{
		Distance: 10,
		Lat:      lat,
		Lon:      lon,
	}

	localEventsQueryTemplate.Execute(writer, data)
	query := elastic.NewRawStringQuery(writer.String())
	return es.client.Search().
		Index(viper.GetString("ES_INDEX")).
		Query(query).
		Pretty(true).
		Do(es.ctx)
}

func (es *ElasticSearch) RollbackES(elasticId string) error {
	_, err := es.client.Delete().
		Index(viper.GetString("ES_INDEX")).
		Type(viper.GetString("ES_INDEX_TYPE")).
		Id(elasticId).
		Do(es.ctx)
	return err
}

// TODO
func (es *ElasticSearch) RollbackUpdateES(elasticId string) error {
	return nil
}

func (es *ElasticSearch) SaveEvent(event *event.Event) (string, error) {

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	resp, err := es.client.Index().
		Index(viper.GetString("ES_INDEX")).
		Type(viper.GetString("ES_INDEX_TYPE")).
		Id(event.EventDetail.Id).
		BodyString(string(eventBytes)).
		Do(es.ctx)

	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (es *ElasticSearch) UpdateEventById (even *event.UpdateRequest) error {

	resp, err :=  es.client.Update().
		Index(viper.GetString("ES_INDEX")).
		Type(viper.GetString( "ES_INDEX_TYPE")).
		Id(even.EventId).
		DetectNoop(true).
		Doc(even.Event).
		Do(es.ctx)

	log.Info("elastic update result:", resp.Result)
	return err
}

func getClient() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL("http://"+viper.GetString("ES_HOST")+":"+viper.GetString("ES_PORT")),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New()),
		elastic.SetInfoLog(log.New()),
		elastic.SetHealthcheck(false),
	)
}

func isPastEvent(event SearchItem) (bool) {

	t, err := time.Parse("3:04 PM", event.EventDetail.EndTime)

	if err != nil {
		log.Error("failed to parse time")
		log.Errorln("Error ", err)
		return true
	}

	formattedTime := t.Format("15:04:05")

	timeString := fmt.Sprintf("%sT%sZ",  event.EventDetail.EndDate,formattedTime)

	eventEndTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		log.Error("failed to parse the event time")
		log.Errorln("Error ", err)
		return true
	}

	if eventEndTime.Before(time.Now()) {
		return true
	}

	return false
}

func (es *ElasticSearch) PublishEvent(eventId, date string) error {

	data := map[string]interface{}{
		"status":     "PUBLISHED",
		"created_on": date,
	}
	resp, err := es.client.Update().
		Index(viper.GetString("ES_INDEX")).
		Type(viper.GetString("ES_INDEX_TYPE")).
		Id(eventId).
		DetectNoop(true).
		Doc(data).
		Do(es.ctx)

	log.Info("elastic update result:", resp.Result)
	return err
}
