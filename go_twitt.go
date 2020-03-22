package main

import (
	"html/template"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"unicode"
	"html"
	twitterscraper "github.com/n0madic/twitter-scraper"
)


type tweet struct {
	Subject string
	Date time.Time
	News  []string
}
func (selfTweet *tweet) FormattedDate() string {
	layout := "2 Jan 2006 15:04:05"
	return selfTweet.Date.Format(layout)
}

type ownTweetTrends struct {
	TweetTrends []tweet
	LastScan time.Time
}
func (ownTweetTrends *ownTweetTrends)httpTrendHandler(w http.ResponseWriter,r *http.Request) {
	log.Info("request site")
	templateFileName := "template_tweet_trend.html"
	t := template.Must(template.New(templateFileName).ParseFiles(templateFileName))
	err := t.Execute(w,ownTweetTrends)	
	if err != nil {
		log.Error(err)
	}
}

func (ownTweetTrends *ownTweetTrends) saveTweet(tweetSubject string) bool {

	for  _,a := range ownTweetTrends.TweetTrends {
		if a.Subject == tweetSubject {
			return false
		}
	}


	response,err := googleRequest(tweetSubject)
	if err != nil {
		log.Error(err)
	}

	results,err := googleResultParser(response)
	if err != nil {
		log.Error(err)
	}

	newTweet := & tweet {
		Subject: tweetSubject,
		Date: time.Now(),
		News: results,
	}
	ownTweetTrends.TweetTrends = append(ownTweetTrends.TweetTrends,*newTweet)
	return true

}
func (ownTweetTrends *ownTweetTrends) GetTrends() {
	trends, err := twitterscraper.GetTrends()
    if err != nil {
        log.Warn("Twitter API not reachable",err)
	}
	for _,a := range trends {
		if ownTweetTrends.saveTweet(a) == true {
			log.Info("Added new Trend: ",a)
		}
	}
	ownTweetTrends.LastScan = time.Now()
}

func (ownTweetTrends *ownTweetTrends) printTrends() {
	log.Info("Trends:")
	for _,a := range ownTweetTrends.TweetTrends {
		log.Info(a.Subject)
	}
}
func googleRequest(queryString string) (*http.Response, error) {
	queryString = strings.ReplaceAll(queryString, " ","%20")
	queryString = strings.ReplaceAll(queryString, "#","")
	
	searchURL := "https://www.google.com/search?q="+queryString+"&hl=en&tbm=nws"
	baseClient := &http.Client{}
 
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	res, err := baseClient.Do(req)
	log.Debug("query: ",searchURL)

	if err != nil {
		return nil, err
	}
	return res, nil

}

func googleResultParser(response *http.Response) ([]string, error){
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	var results []string
	cnt := 0
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
    band, ok := s.Attr("role")
    if ok {
		if band == "heading" {
			if cnt >= 3 {
				return
			} 
			cnt++
			title,_ := s.Html()
			title = strings.Map(func(r rune) rune {
				if unicode.IsGraphic(r) {
					return r
				}
				return -1
			}, title)
			title = strings.Trim(title, "...")
			title = html.UnescapeString(title)
			results = append(results, title)
		}
    }
	})
		log.Debug("found: ",cnt)

	return results, err
}

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.Debug("Logging started")
	reScanTime := 30
	trendStorage := new(ownTweetTrends)
	ticker := time.NewTicker(time.Duration(reScanTime) * time.Second)
	defer ticker.Stop()
	trendStorage.GetTrends()
	http.HandleFunc("/trends",trendStorage.httpTrendHandler)
	
	go http.ListenAndServe(":8090",nil)
	for {
		select {
		case t := <-ticker.C:
			trendStorage.GetTrends()
			log.Info("Scan Trends: ", t)
		}
	}
}