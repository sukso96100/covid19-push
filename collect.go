package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type LastData struct {
	Confirmed int
	Cured     int
	Death     int
	UpdatedAt time.Time
}

type LastNews struct {
	PostId     int
	Title      string
	Department string
	UpdatedAt  time.Time
}

var lData LastData = LastData{}
var lNews LastNews = LastNews{}

statTemplate := "{'confirmed':%d, 'confirmedDiff':%d, 'cured':%d, 'curedDiff':%d, 'death':%d, 'deathDiff':%d}"
newsTemplate := "{'postId':%d, title: '%s', 'department':'%s'}"

func Collect(w http.ResponseWriter, r *http.Request) {
	if lData == (LastData{}) {
		data := GetLastStat()
		lData.Confirmed = data.Confirmed
		lData.Cured = data.Cured
		lData.Death = data.Death
		lData.UpdatedAt = data.UpdatedAt
	}
	if lNews == (LastNews{}) {
		data := GetLastNews()
		lNews.PostId = data.PostId
		lNews.Title = data.Title
		lNews.Department = data.Department
		lNews.UpdatedAt = data.UpdatedAt
	}
	collectData()
	w.WriteHeader(http.StatusOK)
}

func collectData() {
	if lData.UpdatedAt.Add(time.Minute * 1).Before(time.Now()) {
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
		var current = LastData{}
		doc.Find("div.co_cur > ul > li").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			count, _ := strconv.Atoi(s.Find("a").Text())
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
			// save and notify updates
			newData := StatData{
				Confirmed: current.Confirmed,
				Cured:     current.Cured,
				Death:     current.Death,
			}

			newData.Create()
			newStatData := fmt.Sprintf(statTemplate, 
				current.Confirmed, current.Confirmed - lData.Confirmed,
				current.Cured, current.Confirmed - lData.Cured,
				current.Death, current.Death - lData.Death,
			)
			Pusher.Publish("stat", &sse.Event{
				Data: []byte(newStatData),
			})
		}

		lData.UpdatedAt = time.Now()
	}
}
