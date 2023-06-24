package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/podozzoa/couponcrawler/store"
)

var m sync.Mutex

func InitAPI() {
	m.Lock()
	defer m.Unlock()

	router := gin.Default()
	router.GET("/latest", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"title": store.LatestPost.Title, "author": store.LatestPost.Author, "link": store.LatestPost.Link})
	})

	router.Run(":8080")
}
