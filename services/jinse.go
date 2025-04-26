package services

import (
	"encoding/json"
	"io"
	"net/http"
)

var jinseAPI = "https://api.jinse.cn/noah/v2/lives?reading=false&_source=m&flag=down&id=0&category=0&limit=5"

type ListData struct {
	LiveData []LiveData `json:"list"`
}

type LiveData struct {
	Lives []LiveItem `json:"lives"`
}

type LiveItem struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Link      string `json:"link"`
	CreatedAt int64  `json:"created_at"`
}

func GetJinseStories() ([]LiveItem, error) {
	content, _ := getJSON(jinseAPI)
	var f ListData
	if err := json.Unmarshal(content, &f); err != nil {
		return nil, err
	}
	if len(f.LiveData) == 0 {
		return []LiveItem{}, nil
	}
	return f.LiveData[0].Lives, nil
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
	content, err := io.ReadAll(response.Body)
	return content, err
}
