package model

import (
	"encoding/json"
	"fmt"
	"log"
)

type PostData struct {
	From         string `json:"from"`
	Num          int    `json:"num"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	Link         string `json:"link"`
	Crawlingdate string `json:"Crawlingdate"`
}

func (p *PostData) ToMap() (map[string]interface{}, error) {
	jsonPost, err := json.Marshal(p)
	if err != nil {
		log.Printf("Failed to marshal post data: %v", err)
		return nil, err
	}

	var postMap map[string]interface{}
	err = json.Unmarshal(jsonPost, &postMap)
	if err != nil {
		log.Printf("Failed to unmarshal post data: %v", err)
		return nil, err
	}

	return postMap, nil
}

func CreatePostID(postNum int) string {
	return fmt.Sprintf("Post%d", postNum)
}
