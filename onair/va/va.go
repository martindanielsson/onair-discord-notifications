package va

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	baseEndpoint = "https://server1.onair.company/api/v1/"
	userAgent    = "onair-discord-notifications"
)

type Client struct {
	apiKey    string
	baseURL   *url.URL
	client    *http.Client
	common    service
	userAgent string
}

type service struct {
	client *Client
}

func New(httpClient *http.Client, apiKey string) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(baseEndpoint)
	c := &Client{apiKey: apiKey, baseURL: baseURL, client: httpClient, userAgent: userAgent}
	c.common.client = c

	return c, nil
}

func (c *Client) newRequest(endpoint string) (*http.Request, error) {
	u, err := c.baseURL.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	req, err := http.NewRequest("GET", u.String(), buf)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("oa-apikey", c.apiKey)

	return req, nil
}

type response struct {
	Content json.RawMessage `json:"Content"`
}

type Notification struct {
	ID             string  `json:"Id"`
	PeopleID       string  `json:"PeopleId"`
	CompanyID      string  `json:"CompanyId"`
	IsRead         bool    `json:"IsRead"`
	IsNotification bool    `json:"IsNotification"`
	ZuluEventTime  string  `json:"ZuluEventTime"`
	Category       int     `json:"Category"`
	Action         int     `json:"Action"`
	Description    string  `json:"Description"`
	Amount         float64 `json:"Amount"`
	AccountId      string  `json:"AccountId"`
}

type Flight struct {
	ID                   string  `json:"Id"`
	DepartureAirport     Airport `json:"DepartureAirport"`
	ArrivalActualAirport Airport `json:"ArrivalActualAirport"`
	Registered           bool    `json:"Registered"`
	ResultComments       string  `json:"ResultComments"`
	StartTime            string  `json:"StartTime"`
	EndTime              *string `json:"EndTime"`
}

type Airport struct {
	ID   string `json:"Id"`
	ICAO string `json:"ICAO"`
}

type CashFlow struct {
	Entries []CashFlowEntry `json:"Entries"`
}

type CashFlowEntry struct {
	ID           string  `json:"Id"`
	CompanyID    string  `json:"CompanyId"`
	AccountID    string  `json:"AccountId"`
	Amount       float64 `json:"Amount"`
	CreationDate string  `json:"CreationDate"`
	Description  string  `json:"Description"`
	CarryForward string  `json:"CarryFowarad"`
}

func (c *Client) Notifications(vaID string) ([]Notification, *http.Response, error) {
	var r response
	u := "company/" + vaID + "/notifications"

	req, err := c.newRequest(u)

	if err != nil {
		return nil, nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, resp, err
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&r)

	if err != nil {
		return nil, resp, fmt.Errorf("error decoding: %v", err)
	}

	var v []Notification
	err = json.Unmarshal(r.Content, &v)

	return v, resp, err
}

func (c *Client) Flights(vaID string) ([]Flight, *http.Response, error) {
	var r response
	u := "company/" + vaID + "/flights"

	req, err := c.newRequest(u)

	if err != nil {
		return nil, nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, resp, err
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&r)

	if err != nil {
		return nil, resp, fmt.Errorf("error decoding: %v", err)
	}

	var v []Flight
	err = json.Unmarshal(r.Content, &v)

	return v, resp, err
}

func (c *Client) CashFlow(vaID string) (*CashFlow, *http.Response, error) {
	var r response
	u := "company/" + vaID + "/cashflow"

	req, err := c.newRequest(u)

	if err != nil {
		return nil, nil, err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, resp, err
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&r)

	if err != nil {
		return nil, resp, fmt.Errorf("error decoding: %v", err)
	}

	var v CashFlow
	err = json.Unmarshal(r.Content, &v)

	return &v, resp, err
}
