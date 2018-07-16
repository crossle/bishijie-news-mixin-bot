package services

import (
	"encoding/json"
)

var jinseAPI = "https://api.jinse.com/v4/live/list?reading=false&_source=m&limit=5"

type ListData struct {
	LiveData []LiveData `json:"list"`
}

type LiveData struct {
	Lives []LiveItem `json:"lives"`
}

type LiveItem struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

func GetJinseStories() ([]LiveItem, error) {
	content, _ := getJSON(jinseAPI)
	var f ListData
	if err := json.Unmarshal(content, &f); err != nil {
		return nil, err
	}
	return f.LiveData[0].Lives, nil
}
