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
	"github.com/xuri/excelize/v2"
)

// 1.http + robotstxt + excel

// Information
type Information struct {
	InfoID   int
	InfoLink string
	Time     int64
}

var UserAgent string = "/"

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
func ParseHtml(htmlcontent string) (bool, []Information) {
	regA := regexp.MustCompile(`href="(.*?)"`)
	htmlcontentB := regA.FindAllStringSubmatch(htmlcontent, -1)

	var InfoGroup []Information
	for k, v := range htmlcontentB {
		var Info Information
		Info.InfoID = k
		Info.InfoLink = v[1]
		Info.Time = time.Now().UTC().Unix()

		InfoGroup = append(InfoGroup, Info)
		fmt.Println(Info)
		fmt.Println("--- --- ---")
		if k == 10 {
			break
		}
	}
	return true, InfoGroup
}

// ExcelInit
func ExcelInit() {
	f := excelize.NewFile()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet2", "A1", "InfoID")
	f.SetCellValue("Sheet2", "B1", "InfoLink")
	f.SetCellValue("Sheet2", "C1", "Time")
	f.SetActiveSheet(index)
	if err := f.SaveAs("Data.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// ExcelHandler
func ExcelHandler(information []Information) {
	f, err := excelize.OpenFile("Data.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(information); i++ {
		f.SetCellValue("Sheet2", fmt.Sprintf("A%d", i+2), information[i].InfoID)
		f.SetCellValue("Sheet2", fmt.Sprintf("B%d", i+2), information[i].InfoLink)
		f.SetCellValue("Sheet2", fmt.Sprintf("C%d", i+2), information[i].Time)
	}

	if err := f.SaveAs("Data.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func main() {
	CrawlUrl := ""
	BaseUrl := ""

	flag := PreCheckRobotsTxt(CrawlUrl, BaseUrl)
	if flag {
		htmlcontent := Fetch(CrawlUrl)
		flag, infoGroup := ParseHtml(htmlcontent)
		if flag == true {
			ExcelInit()
			ExcelHandler(infoGroup)
		}
	} else {
		fmt.Println("Not Allow")
	}
}
