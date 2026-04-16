package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	feedURL          = "https://javiersilvestri.vercel.app/rss.xml"
	templateFile     = "README.md.tpl"
	outputFile       = "README.md"
	placeholder      = "%{{latest_articles}}%"
	numberOfArticles = 5
)

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

func main() {
	resp, err := http.Get(feedURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error fetching RSS feed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		fmt.Fprintln(os.Stderr, "error decoding XML:", err)
		os.Exit(1)
	}

	items := rss.Channel.Items
	count := min(numberOfArticles, len(items))

	var sb strings.Builder
	for _, item := range items[:count] {
		sb.WriteString(fmt.Sprintf(
			"### [%s](%s)\n<kbd>📅 %s</kbd>\n\n> %s\n\n<br/>\n\n",
			item.Title,
			item.Link,
			formatDate(item.PubDate),
			item.Description,
		))
	}

	tmpl, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading template:", err)
		os.Exit(1)
	}

	output := strings.Replace(string(tmpl), placeholder, sb.String(), 1)

	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "error writing README.md:", err)
		os.Exit(1)
	}

	fmt.Printf("README.md updated (%d posts)\n", count)
}

func formatDate(date string) string {
	parsed, err := time.Parse(time.RFC1123, date)
	if err != nil {
		parsed, err = time.Parse(time.RFC1123Z, date)
		if err != nil {
			return date
		}
	}
	return parsed.Format("January 2, 2006")
}

