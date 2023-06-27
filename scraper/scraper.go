package scraper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/podozzoa/couponcrawler/model"
	"github.com/podozzoa/couponcrawler/store"
)

func CheckNewPosts(ctx context.Context) {
	url := "https://www.inven.co.kr/board/fifaonline4/3145?category=%EC%BF%A0%ED%8F%B0"
	c := colly.NewCollector()

	postList := []model.PostData{}
	from := "피파인벤"

	c.OnHTML("form[name=board_list1] tbody tr:not(.notice)", func(e *colly.HTMLElement) {
		num, _ := strconv.Atoi(e.ChildText("td.num span"))
		title := filteringStrings(e.ChildText("td.tit a.subject-link"))
		author := e.ChildText("td.user span")
		link := e.Request.AbsoluteURL(e.ChildAttr("td.tit a.subject-link", "href"))

		post := model.PostData{From: from, Num: num, Title: title, Author: author, Link: link, Crawlingdate: time.Now().Format("2006-01-02T15:04:05")}
		postList = append(postList, post)
	})

	c.Visit(url)

	store.SavePosts(ctx, postList)
}

func filteringStrings(str string) string {
	str = strings.Replace(str, ",", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "[쿠폰]", "", -1)
	str = strings.Replace(str, "!", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}
