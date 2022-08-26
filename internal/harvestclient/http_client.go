package harvestclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURLHarvest = "https://api.harvestapp.com/v2"

func (args Arguments) ToURLValues() url.Values {
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

func (c HarvestClient) get(path string, args Arguments, target interface{}) error {
	url := BaseURLHarvest + path
	urlWithParams := fmt.Sprintf("%s?%s", url, args.ToURLValues().Encode())

	req, err := http.NewRequest("GET", urlWithParams, nil)

	if err != nil {
		return err
	}

	c.AddRequestHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &target)
	if err != nil {
		return err
	}

	return nil
}
