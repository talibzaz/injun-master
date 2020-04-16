package analytics

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	methodGet = "GET"
)

type analyticsClient struct {
	apiKey 	string
	secretKey 	string
}

func newAnalyticsAuthClient() *analyticsClient{
	return &analyticsClient{
		apiKey: "def6005810e873f9bb90521922ad50d1",
		secretKey: "759d6b05d8430bf4357654297bd75645",
	}
}

type clientAmplitudeId struct {
	Matches []struct {
		UserID      string `json:"user_id"`
		AmplitudeID int64  `json:"amplitude_id"`
	} `json:"matches"`
}

type eventDetails struct {
	UserData struct {
		NumEvents            int           `json:"num_events"`
	} `json:"userData"`
}


func (c *analyticsClient) amplitudeRequest(req *http.Request) ([]byte, error){
	req.SetBasicAuth(c.apiKey, c.secretKey)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if 200 != res.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}


//Page Views for single event
func GetPageViewsByEventID(eventId string) (int64, error) {
	client := newAnalyticsAuthClient()

	res, err := getAmplitudeId(eventId,client)
	if err != nil {
		return 0, err
	}

	amplitudeID := strconv.FormatInt(res, 10)

	eventHits, err := getEventPageHits(amplitudeID, client)
	if err != nil {
		return 0, err
	}

	return int64(eventHits), nil
}

func getAmplitudeId(eventId string, client *analyticsClient) (int64, error){
	URL := `https://amplitude.com/api/2/usersearch?user=` + eventId

	req, err := http.NewRequest(methodGet, URL, nil)
	if err != nil {
		return 0, err
	}
	body, err := client.amplitudeRequest(req)
	if err != nil {
		return 0, err
	}

	var clientData clientAmplitudeId

	err = json.Unmarshal(body, &clientData)
	if err != nil {
		return 0, err
	}
	if len(clientData.Matches)  > 0 {
		return clientData.Matches[0].AmplitudeID, nil

	}
	return 0, nil
}

func getEventPageHits(amplitudeID string, client *analyticsClient) (int,error) {
	URL := `https://amplitude.com/api/2/useractivity?user=` + amplitudeID

	req, err := http.NewRequest(methodGet, URL, nil)
	if err != nil {
		return 0, err
	}

	body, err := client.amplitudeRequest(req)
	if err != nil {
		return 0, err
	}

	var eventHits eventDetails

	err = json.Unmarshal(body, &eventHits)
	if err != nil {
		return 0, err
	}

	return eventHits.UserData.NumEvents, nil
}