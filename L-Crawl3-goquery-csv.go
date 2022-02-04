package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"encoding/csv"

	"github.com/PuerkitoBio/goquery"
)

// 3.2 net/http goquery  csv

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

func parseUrls(url string, ch chan bool, w *csv.Writer) {
	doc := fetch(url)
	doc.Find(".media-item").Each(func(index int, s *goquery.Selection) {
		// for each item
		title := s.Find("H5").Text()
		fmt.Println(index, title)

		fmt.Println("--- --- ---")
		err := w.Write([]string{strconv.Itoa(index), title})
		checkError(err)
	})
	ch <- true
}

func main() {
	start := time.Now()
	url := ""
	ch := make(chan bool)

	// save csv
	f, err := os.Create("abc.csv")
	checkError(err)
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.Write([]string{"ID", "Title"})
	checkError(err)

	for i := 1; i <= 5; i++ {
		// newurl := fmt.Sprint(url, i)
		go parseUrls(url+strconv.Itoa(i), ch, writer)
	}

	for i := 1; i <= 5; i++ {
		<-ch
	}
	f.Sync()

	elapsed := time.Since(start)
	fmt.Println("Took %s", elapsed)
}
