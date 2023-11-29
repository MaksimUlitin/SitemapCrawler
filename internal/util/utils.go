package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/maksimulitin/pkg/model"
)

func RandomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(model.UserAgents)
	return model.UserAgents[randNum]
}

func IsSitemap(urls []string) ([]string, []string) {
	sitemapFiles := []string{}
	pages := []string{}
	for _, page := range urls {
		foundSitemap := strings.Contains(page, "xml")
		if foundSitemap == true {
			fmt.Println("Found Sitemap", page)
			sitemapFiles = append(sitemapFiles, page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFiles, pages
}
