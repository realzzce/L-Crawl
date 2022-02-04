package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/temoto/robotstxt"
)

var UserAgent string = "-"

// fetch
func Fetch(url string) string {
	fmt.Println("url: ", url)

	// request url
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UserAgent)

	// create http client if x509 err
	timeout := time.Duration(10 * time.Second)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	// send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println("Http status code: ", resp.StatusCode)
	}

	// contents body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error: ", err)
		return ""
	}

	return string(body)
}

// PreCheckRobotsTxt
func PreCheckRobotsTxt(CrawlUrl, BaseUrl string) bool {
	robotstxturl := ""
	if BaseUrl == "" {
		reg := regexp.MustCompile("^(http.*?//.*?com)/")
		hosturl := reg.FindAllStringSubmatch(CrawlUrl, -1)
		BaseUrl = hosturl[0][1]
	}

	robotstxturl = BaseUrl + "/robots.txt"

	robotstxtcontents := Fetch(robotstxturl)
	robots, err := robotstxt.FromString(robotstxtcontents)
	if err != nil {
		fmt.Println(err)
	}

	checkUrl := strings.Split(CrawlUrl, BaseUrl)[1]
	if checkUrl == "" {
		checkUrl = "/"
	}
	allowed := robots.TestAgent(checkUrl, UserAgent)

	return allowed
}

// ParseHtml
func ParseHtml(htmlcontent string) bool {
	regA := regexp.MustCompile(`href="(.*?)"`)
	htmlcontentB := regA.FindAllStringSubmatch(htmlcontent, -1)

	var InfoGroup []interface{}
	for k, v := range htmlcontentB {
		Info := make(map[string]interface{})
		Info["InfoID"] = k
		Info["InfoLink"] = v[1]
		Info["Time"] = time.Now().UTC().Unix()

		InfoGroup = append(InfoGroup, Info)
		fmt.Println(Info)
		fmt.Println("--- --- ---")
		if k == 10 {
			break
		}
	}
	fmt.Println(InfoGroup)
	return true
}

func main() {
	CrawlUrl := "https://oohtalk.com"
	// BaseUrl := "https://oohtalk.com"

	// flag := PreCheckRobotsTxt(CrawlUrl, BaseUrl)
	if true {
		htmlcontent := Fetch(CrawlUrl)
		fmt.Println(htmlcontent)
	} else {
		fmt.Println("Not Allow")
	}
}
