package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/podozzoa/couponcrawler/model"
)

var m sync.Mutex
var latestPost model.PostData

func InitAPI() {
	m.Lock()
	defer m.Unlock()

	router := gin.Default()
	router.GET("/latest", func(c *gin.Context) {
		m.Lock()
		defer m.Unlock()
		c.JSON(http.StatusOK, gin.H{"title": latestPost.Title, "author": latestPost.Author, "link": latestPost.Link})
	})

	router.Run(":8080")
}
