package salmongo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const subscriptionsURL string = "/redfish/v1/EventService/Subscriptions/"

type Client struct {
	HostFQDN string
	Username string
	Password string
}

type Subscription struct {
	Context     string   `json:"Context"`
	Destination string   `json:"Destination"`
	EventTypes  []string `json:"EventTypes"`
	Protocol    string   `json:"Protocol"`
	Id          string   `json:"Id,omitempty"`
}

type RemovalResponse struct {
	Message interface{} `json:"@Message.ExtendedInfo"`
}

func SalmonClient(hostfqdn string, username string, password string) *Client {
	return &Client{
		HostFQDN: hostfqdn,
		Username: username,
		Password: password,
	}
}

func (s *Client) CreateSubscription(subscription *Subscription) (*Subscription, error) {
	url := fmt.Sprintf("https://%s"+subscriptionsURL, s.HostFQDN)
	jsondata, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
	req.Header.Set("Content-Type", "application/json")
	bytes, err := s.sendRequest(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	var data Subscription
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Client) GetSubscription(uuid string) (*Subscription, error) {
	url := fmt.Sprintf("https://%s"+subscriptionsURL, s.HostFQDN, uuid)
	req, _ := http.NewRequest("GET", url, nil)
	bytes, err := s.sendRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var data Subscription
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Client) DeleteSubscription(uuid string) (*RemovalResponse, error) {
	url := fmt.Sprintf("https://%s"+subscriptionsURL+"%s", s.HostFQDN, uuid)
	req, _ := http.NewRequest("DELETE", url, nil)
	bytes, err := s.sendRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var data RemovalResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Client) sendRequest(req *http.Request, wantedResponse int) ([]byte, error) {
	req.SetBasicAuth(s.Username, s.Password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Error: Got HTTP 401, Unauthorized")
	}
	if response.StatusCode != wantedResponse {
		return nil, fmt.Errorf(string(data))
	}
	return data, nil
}
