package main

import (
  "encoding/xml"
  "fmt"
  "time"
  "net/http"
  "os"
  "strings"
)

const (
  NUMBER_OF_ARTICLES = 3
)

type Item struct {
  Title       string `xml:"title"`
  Link        string `xml:"link"`
  Description string `xml:"description"`
  PubDate     string `xml:"pubDate"`
}

type Channel struct {
  Title       string `xml:"title"`
  Description string `xml:"description"`
  Link        string `xml:"link"`
  Items       []Item `xml:"item"`
}

type RSS struct {
  Channel Channel `xml:"channel"`
}

func main() {
  // Fetch RSS feed
  resp, err := http.Get("https://javiersilvestri.vercel.app/rss.xml")
  if err != nil {
    fmt.Println("Error fetching RSS feed:", err)
    return
  }
  defer resp.Body.Close()

  // Parse XML
  rss := RSS{}
  err = xml.NewDecoder(resp.Body).Decode(&rss)
  if err != nil {
    fmt.Println("Error decoding XML:", err)
    return
  }

  // Prepare markdown content for latest articles
  var latestArticlesMarkdown strings.Builder
  for _, item := range rss.Channel.Items[:NUMBER_OF_ARTICLES] {
    // latestArticlesMarkdown.WriteString(fmt.Sprintf("- [%s](%s)\n", item.Title, item.Link))
    latestArticlesMarkdown.WriteString(fmt.Sprintf(`
  - [%s](%s)

    ► Published on: %s
    - %s
 ---
    `,
    item.Title,
    item.Link,
    formatDate(item.PubDate),
    item.Description,
  ))
  }

  // Read the README template
  templateContent, err := os.ReadFile("README.md.tpl")
  if err != nil {
    fmt.Println("Error reading README.md.tpl:", err)
    return
  }

  // Replace placeholders in the template
  newMarkdown := strings.Replace(string(templateContent), "%{{latest_articles}}%", latestArticlesMarkdown.String(), 1)

  // Write to README.md
  err = os.WriteFile("README.md", []byte(newMarkdown), os.ModePerm)
  if err != nil {
    fmt.Println("Error writing to README.md:", err)
    return
  }

  fmt.Println("README.md file updated successfully")
}

func formatDate(date string) string {
  parsedDate, err := time.Parse(time.RFC1123, date)
  if err != nil {
    return date
  }
  return parsedDate.Format("January 2, 2006")
}
