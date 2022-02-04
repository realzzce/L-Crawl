package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// 5. htmlquery (xpath)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func fetch(url string) *html.Node {
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

	doc, err := htmlquery.Parse(res.Body)
	checkError(err)

	return doc
}

func parseUrls(url string) {
	doc := fetch(url)
	nodes := htmlquery.Find(doc, `//*[@id="media-box"]/div[*]/div[2]`)
	for _, node := range nodes {
		//*[@id="media-box"]/div[1]/div[2]/a
		//*[@id="media-box"]/div[1]/div[2]/a/h5

		// links := htmlquery.FindOne(node, "//a[@href]")
		// link := htmlquery.SelectAttr(links, "href")

		links := htmlquery.FindOne(node, "./a/@href")
		link := htmlquery.InnerText(links)

		titles := htmlquery.FindOne(node, `./a/h5/text()`)
		title := htmlquery.InnerText(titles)

		fmt.Println(link)
		fmt.Println(title)
		fmt.Println("--- --- ---")
	}
}

func main() {
	url := ""
	for i := 1; i <= 5; i++ {
		parseUrls(url + strconv.Itoa(i))
	}
}
