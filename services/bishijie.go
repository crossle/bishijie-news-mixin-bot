package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var apiEndpoint = "https://api.bishijie.com/news"

type HoldData struct {
	Data map[string]Buttom `json:"data"`
}

type Buttom struct {
	Items []NewsFlash `json:"buttom"`
}

type NewsFlash struct {
	ID      int    `json:"newsflash_id"`
	Content string `json:"content"`
}

func GetStories() ([]NewsFlash, error) {
	content, _ := getJSON(apiEndpoint)
	var f HoldData
	if err := json.Unmarshal(content, &f); err != nil {
		return nil, err
	}
	for _, v := range f.Data {
		return v.Items, nil
	}
	return nil, nil
}

func getJSON(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Close = true // connection reset

	client := new(http.Client)
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	return content, err
}
