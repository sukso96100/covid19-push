package main

import (
	"context"
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

	"github.com/chromedp/chromedp"
)

func Collect(c echo.Context) error {
	collectStat()
	// collectNews()
	return c.String(http.StatusOK, "OK")
}

func collectStat() {
	lastStat := database.GetLastStat()
	if lastStat.CreatedAt.Add(time.Second * 10).Before(time.Now()) {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		contextVar, cancelFunc := chromedp.NewContext(
			allocCtx,
			chromedp.WithLogf(log.Printf),
		)
		defer cancelFunc()
		var content string
		err := chromedp.Run(contextVar,
			chromedp.Navigate(`http://ncov.mohw.go.kr/bdBoardList_Real.do?brdId=1&brdGubun=11`),
			chromedp.WaitVisible(`body`, chromedp.ByQuery),
			chromedp.Sleep(time.Second*3),
			chromedp.InnerHTML(`body`, &content, chromedp.ByQuery),
		)
		if err != nil {
			fmt.Println("%w", err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			log.Fatal(err)
		}

		items := doc.Find("table").Eq(0).Find("tr")
		current := database.StatData{
			Confirmed: convertToInt(items.Eq(0).Find("td").Text()),
			Cured:     convertToInt(items.Eq(1).Find("td").Text()),
			Death:     convertToInt(items.Eq(2).Find("td").Text()),
			Checking:  convertToInt(items.Eq(3).Find("td").Text()),
		}
		fmt.Println(items.Eq(0).Find("td").Text())
		fmt.Println(convertToInt(items.Eq(0).Find("td").Text()))
		if lastStat.Confirmed != current.Confirmed ||
			lastStat.Cured != current.Cured ||
			lastStat.Death != current.Death ||
			lastStat.Checking != current.Checking {
			fmt.Println("Sending new stat data")
			current.Create()
			fcm.GetFCMApp().PushStatData(
				lastStat, current,
			)
		}
	}
}

func convertToInt(str string) int {
	r1 := strings.Split(str, "명")[0]
	r2 := strings.ReplaceAll(r1, ",", "")
	result, _ := strconv.Atoi(strings.TrimSpace(r2))
	fmt.Println(r1, r2, result)
	return result
}

func collectNews() {
	var lNews database.NewsData = database.GetLastNews()
	if lNews.UpdatedAt.Add(time.Second * 1).Before(time.Now()) {
		fmt.Println("Collecting news data...")
		// collect data
		// Request the HTML page.
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req, err := http.NewRequest("GET", "http://ncov.mohw.go.kr/tcmBoardList.do?brdId=3", nil)
		req.Close = true
		res, err := client.Do(req)
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
