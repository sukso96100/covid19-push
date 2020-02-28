package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/sukso96100/covid19-push/database"
	"github.com/sukso96100/covid19-push/fcm"
	"github.com/labstack/echo/v4"

	// "io/ioutil"
	"strings"
)

var lData database.StatData = database.StatData{}
var lNews database.NewsData = database.NewsData{}

const statTemplate = "{'confirmed':%d, 'confirmedDiff':%d, 'cured':%d, 'curedDiff':%d, 'death':%d, 'deathDiff':%d}"
const newsTemplate = "{'postId':%d, title: '%s', 'department':'%s'}"

func Collect(c echo.Context) error { {
	collectData()
	
	return c.String(http.StatusOK, "OK")
}

func collectData() {
	if lData.UpdatedAt.Add(time.Second * 1).Before(time.Now()) {
		fmt.Println("Collecting data...")
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
		var current = database.StatData{}
		doc.Find("div.co_cur > ul > li").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			raw := s.Find("a").Text()
			count, _ := strconv.Atoi(strings.Split(raw, " ")[0])
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
				current.Confirmed-lData.Cured,
				current.Death-lData.Death,
			)

		}

		lData = current
	}
}
