package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/sukso96100/covid19-push/database"
	"github.com/sukso96100/covid19-push/fcm"

	// "io/ioutil"
	"strings"
)

func Collect(c echo.Context) error {
	collectStat()
	collectNews()
	return c.String(http.StatusOK, "OK")
}

func collectStat() {
	var lData database.StatData = database.GetLastStat()
	if lData.UpdatedAt.Add(time.Second * 1).Before(time.Now()) {
		fmt.Println("Collecting stat data...")
		// collect data
		// Request the HTML page.
		res, err := http.Get("http://ncov.mohw.go.kr/index_main.jsp")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		var current = database.StatData{
            Confirmed: lData.Confirmed,
            Cured: lData.Cured,
            Death: lData.Death,
        }
		doc.Find("div.co_cur > ul > li").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			raw := s.Find("a").Text()
			fmt.Println(raw)
			count, _ := strconv.Atoi(strings.ReplaceAll(strings.Split(raw, " ")[0], ",", ""))
			fmt.Println(count)
			switch i {
			case 0:
				current.Confirmed = count
			case 1:
				current.Cured = count
			case 2:
				current.Death = count
			}
		})
		if lData.Confirmed != current.Confirmed ||
			lData.Cured != current.Cured ||
			lData.Death != current.Death {
			fmt.Println("Notifying stat updates...")
			// save and notify updates
			current.Create()
			fcm.GetFCMApp().PushStatData(
				current,
				current.Confirmed-lData.Confirmed,
				current.Cured-lData.Cured,
				current.Death-lData.Death,
			)

		}
	}
}

func collectNews() {
    newsUrl := "http://ncov.mohw.go.kr/tcmBoardList.do?brdId=&brdGubun=&dataGubun=&ncvContSeq=&contSeq=&board_id="
	var lNews database.NewsData = database.GetLastNews()
	if lNews.UpdatedAt.Add(time.Second * 1).Before(time.Now()) {
		fmt.Println("Collecting news data...")
		// collect data
		// Request the HTML page.
		res, err := http.Get(newsUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fn_tcm_boardView('/tcmBoardView.do','','','353254','', 'ALL');
		// http://ncov.mohw.go.kr/tcmBoardView.do?ncvContSeq=353254&contSeq=353254&gubun=ALL
		tds := doc.Find("tbody > tr").First().Find("td")

		linkFunc := tds.Eq(1).Find("a").AttrOr("onclick", "")
		newsLink := "http://ncov.mohw.go.kr/tcmBoardList.do"
		if linkFunc != "" {
			tmpl := "http://ncov.mohw.go.kr/tcmBoardView.do?ncvContSeq=%s&contSeq=%s&gubun=ALL"
			splits := strings.Split(linkFunc, ",")
			postNum := strings.ReplaceAll(splits[3], "'", "")
			newsLink = fmt.Sprintf(tmpl, postNum, postNum)
		}
		postNum, _ := strconv.Atoi(tds.Eq(0).Text())
		current := database.NewsData{
			PostId:     postNum,
			Title:      tds.Eq(1).Find("a").Text(),
			Department: tds.Eq(2).Text(),
			Link:       newsLink,
		}
		fmt.Println(current.Title)
		if lNews.PostId != current.PostId {
			fmt.Println("Notifying news updates...")
			current.Create()
			fcm.GetFCMApp().PushNewsData(current)
		}
	}
}
