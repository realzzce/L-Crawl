package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// 2.chromedp

// chromedpFetch chromedp run
func chromedpFetch(url string) string {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false), // headless
		chromedp.UserAgent("-"),          // user-agent
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// create chrome init
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// run
	var resHtml string
	err := chromedp.Run(ctx,
		visiteURL(url, &resHtml),
	)
	if err != nil {
		log.Fatal(err)
	}

	return resHtml
}

// visiteURL get outhtml
func visiteURL(url string, resHtml *string) chromedp.Tasks {
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`),
		chromedp.Sleep(3 * time.Second),
		chromedp.Click(".nav-login-btn", chromedp.BySearch),
		chromedp.Sleep(3 * time.Second),
		chromedp.OuterHTML("html", resHtml, chromedp.ByQuery),
		// chromedp.OuterHTML(".talk-title-index-info", &body, chromedp.ByQuery),
		chromedp.Sleep(3 * time.Second),
		chromedp.SendKeys(`.nav-search-input`, "abc"),
		chromedp.SendKeys(`.nav-search-input`, kb.Enter),
		chromedp.Sleep(3 * time.Second),
		chromedp.EvaluateAsDevTools(`alert('aaaaa');`, nil),
		chromedp.Sleep(3 * time.Second),
		chromedp.EvaluateAsDevTools(`alert('bbbb');`, nil),
		chromedp.Sleep(3 * time.Second),
	}
	return tasks
}

func main() {
	var url string
	url = ``
	// 1. get html
	resHtml := chromedpFetch(url)
	// 2. parse html
	fmt.Println(resHtml)
}
