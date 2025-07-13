package main

import (
	"fmt"

	"github.com/awaisahmedsukhera/sykell-crawler-task/backend/internal/crawler"
)

func main() {
	// Todo: Remove after development and testing
	result, err := crawler.CrawlURL("https://example.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("HTMLVersion:", result.HTMLVersion)
	fmt.Println("PageTitle:", result.PageTitle)
}
