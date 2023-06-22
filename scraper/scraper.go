package scraper

import (
	"context"
	"strconv"

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
		title := e.ChildText("td.tit a.subject-link")
		author := e.ChildText("td.user span")
		link := e.Request.AbsoluteURL(e.ChildAttr("td.tit a.subject-link", "href"))

		post := model.PostData{From: from, Num: num, Title: title, Author: author, Link: link}
		postList = append(postList, post)
	})

	c.Visit(url)

	store.SavePosts(ctx, postList)
}