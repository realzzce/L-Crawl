package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 3.3 net/http goquery mysql

type Topic struct {
	ID    int    `gorm:"AUTO_INCREATEMENT"`
	Title string `gorm:"type:varchar(120);"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// fetch url
func fetch(url string) *goquery.Document {
	fmt.Println(url)
	// x509
	timeout := time.Duration(10 * time.Second)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	// request url
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "-")

	// send request
	res, err := client.Do(req)
	checkError(err)

	// res body
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal("status code error: %d %s", res.StatusCode, res.Status)
	}

	// load html
	fmt.Println("--- --- ---")
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	return doc
}

func parseUrls(url string, ch chan bool, db *gorm.DB) {
	doc := fetch(url)
	doc.Find(".media-item").Each(func(index int, s *goquery.Selection) {
		// for each item
		title := s.Find("H5").Text()
		fmt.Println(index, title)

		fmt.Println("--- --- ---")
		topic := &Topic{
			// ID:    index,
			Title: title,
		}
		db.Create(&topic)

	})
	time.Sleep(2 * time.Second)
	ch <- true
}

func main() {
	start := time.Now()
	url := ""
	ch := make(chan bool)

	// db
	db, err := gorm.Open("mysql", "root:@/test?charset=utf8")
	defer db.Close()
	checkError(err)

	db.DropTableIfExists(&Topic{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Topic{})

	for i := 1; i <= 5; i++ {
		// newurl := fmt.Sprint(url, i)
		go parseUrls(url+strconv.Itoa(i), ch, db)
	}

	for i := 1; i <= 5; i++ {
		<-ch
	}

	elapsed := time.Since(start)
	fmt.Println("Took %s", elapsed)
}
