package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexandrevicenzi/go-sse"

	// "io/ioutil"
	"strings"
)

var lData StatData = StatData{}
var lNews NewsData = NewsData{}

const statTemplate = "{'confirmed':%d, 'confirmedDiff':%d, 'cured':%d, 'curedDiff':%d, 'death':%d, 'deathDiff':%d}"
const newsTemplate = "{'postId':%d, title: '%s', 'department':'%s'}"

func Collect(w http.ResponseWriter, r *http.Request) {
	Pusher.SendMessage("/updates/stat", sse.SimpleMessage("ping"))
	if lData == (StatData{}) {
		data := GetLastStat()
		lData = data
	}
	if lNews == (NewsData{}) {
		data := GetLastNews()
		lNews = data
	}
	collectData()
	w.WriteHeader(http.StatusOK)
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

		// if res.StatusCode == http.StatusOK {
		// 	bodyBytes, err := ioutil.ReadAll(res.Body)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	// bodyString := string(bodyBytes)
		// 	// fmt.Println(bodyString)
		// }

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		var current = StatData{}
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
			newStatData := fmt.Sprintf(statTemplate,
				current.Confirmed, current.Confirmed-lData.Confirmed,
				current.Cured, current.Confirmed-lData.Cured,
				current.Death, current.Death-lData.Death,
			)
			Pusher.SendMessage("/updates/stat", sse.SimpleMessage(newStatData))
		}

		lData = current
	}
}
