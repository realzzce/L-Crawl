package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/anaskhan96/soup"
)

// 4. soup

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// fetch url
func fetch(url string) soup.Root {
	fmt.Println(url)

	// soup.Headers = map[string]string{
	// 	"User-Agent": "-",
	// }
	// source, err := soup.Get(url)
	// checkError(err)

	// create http client if x509 err
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

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	source := buf.String()

	doc := soup.HTMLParse(source)
	return doc

}

func parseUrls(url string) {
	doc := fetch(url)

	html := doc.Find("div", "class", "media-box").FindAll("div", "class", "media-title")
	for _, content := range html {
		link, _ := content.Find("a").Attrs()["href"]
		title := content.Find("h5").Text()
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
