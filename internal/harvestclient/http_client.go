package harvestclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURLHarvest = "https://api.harvestapp.com/v2"

func (args Arguments) ToUrlValues() url.Values {
	v := url.Values{}
	for key, value := range args {
		v.Set(key, value)
	}
	return v
}

func (c HarvestClient) AddRequestHeaders(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.Header.Add("Harvest-Account-Id", c.AccountID)
}

func makeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func unmarshalBody(body []byte, target interface{}) error {
	err := json.Unmarshal(body, &target)

	if err != nil {
		return err
	}

	return nil
}

func makeHttpRequest(req *http.Request, target interface{}) error {
	resp, err := makeRequest(req)
	if err != nil {
		return err
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}

	return unmarshalBody(body, &target)
}

func (c HarvestClient) get(path string, args Arguments, target interface{}) error {
	url := BaseURLHarvest + path
	urlWithParams := fmt.Sprintf("%s?%s", url, args.ToUrlValues().Encode())

	req, err := http.NewRequest("GET", urlWithParams, nil)

	if err != nil {
		return err
	}

	c.AddRequestHeaders(req)

	return makeHttpRequest(req, target)
}
