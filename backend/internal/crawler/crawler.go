package crawler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type CrawlResult struct {
	HTMLVersion string
	PageTitle   string
	// HeadingCount       map[string]int
	// InternalLinksCount int
	// ExternalLinksCount int
	// BrokenLinksCount   int
	// LoginFormPresent   bool
}

func CrawlURL(url string) (*CrawlResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Todo: Remove it after development. This is only for testing
	fmt.Println("===== RAW HTML =====")
	html.Render(os.Stdout, doc)
	fmt.Println("\n===== END =====")

	version := extractHTMLVersion(doc)
	pageTitle := extractPageTitle(doc)

	return &CrawlResult{
		HTMLVersion: version,
		PageTitle:   pageTitle,
	}, nil
}

func extractHTMLVersion(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.DoctypeNode {
			doctype := strings.ToLower(c.Data)

			if doctype == "html" {
				// Likely HTML5
				return "HTML5"
			}

			// Older doctypes may contain more detail in the attributes
			if c.Data == "html" && c.Attr != nil {
				for _, attr := range c.Attr {
					if strings.Contains(attr.Val, "4.01") {
						if strings.Contains(attr.Val, "Transitional") {
							return "HTML 4.01 Transitional"
						}
						if strings.Contains(attr.Val, "Strict") {
							return "HTML 4.01 Strict"
						}
						return "HTML 4.01"
					}
				}
			}

			// Todo: Add other versions if needed

			// Fallback
			return "Unknown DOCTYPE"
		}
	}
	return "DOCTYPE not found"
}

func extractPageTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return strings.TrimSpace(n.FirstChild.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title := extractPageTitle(c)
		if title != "" {
			return title
		}
	}
	return ""
}
