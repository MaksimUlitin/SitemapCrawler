package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type SeoData struct {
	URL             string
	Title           string
	H1              string
	MetaDescription string
	StatusCode      int
}

type DefaultParser struct {
}

type Parser interface {
}

var UserAgentes = []string{}

func main() {
	p := DefaultParser{}
	results := ScrapeSiteMap()
	for _, res := range results {
		fmt.Println(res)
	}
}

func RandomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNam := rand.Int() % len(UserAgentes)
	return UserAgentes[randNam]
}

func isSitemap(urls []string) ([]string, []string) {
	sitemapFile := []string{}
	pages := []string{}

	for _, page := range urls {
		foundSistemFile := strings.Contains(page, "XML")
		if foundSistemFile == true {
			fmt.Println("Found Sistem", page)
			sitemapFile = append(sitemapFile, page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFile, pages
}

func ExtractSiteMapURls(startURL string) []string {
	Worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++
	go func() { Worklist <- []string{startURL} }()
	for ; n > 0; n-- {
	}
	list := <-Worklist
	for _, link := range list {
		n++
		go func(link string) {
			response, err := MakeRequest(link)
			if err != nil {
				log.Printf("error retrieving URL: %s ", link)
			}
			urls, _ := ExtractSiteMapURls(response)
			if err != nil {
				log.Printf("Error extracting document from response, URL: %s", link)
			}
			sitemapFiles, pages := isSitemap(urls)
			if sitemapFiles != nil {
				Worklist <- sitemapFiles
			}
			for _, page := range pages {
				toCrawl = append(toCrawl, page)
			}

		}(link)
	}
}

func MakeRequest() {

}

func ScrapeURLs() {

}

func ScrapePage() {

}

func CrawlPage() {

}

func GetSEOData() {

}

func ScrapeSiteMap(url string) []SeoData {
	result := ExtractSiteMapURls(url)
	res := ScrapeURLs(result)
	return res
}
