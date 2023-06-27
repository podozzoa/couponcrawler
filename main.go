package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/podozzoa/couponcrawler/api"
	"github.com/podozzoa/couponcrawler/scraper"
	"github.com/podozzoa/couponcrawler/store"
)

func main() {
	ctx := context.Background()

	configFile := "config.json"

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	store.InitFirestoreClient(ctx)
	defer store.CloseFirestoreClient()

	go func() {
		for {
			scraper.CheckNewPosts(ctx)
			fmt.Println("Crawling...", time.Now().Format("2006-01-02T15:04:05"))
			time.Sleep(time.Duration(config.CrawlingIntervalSeconds) * time.Second) // 크롤러를 30초마다 실행합니다.
		}
	}()

	api.InitAPI()
}
