package main

import (
	"fmt"

	"github.com/maksimulitin/internal/app"
	"github.com/maksimulitin/pkg/model"
)

func main() {
	p := model.DefaultParser{}

	results := app.ScrapeSitemap("https://", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}
