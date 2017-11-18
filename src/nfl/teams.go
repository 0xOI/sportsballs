package main

import (
  "fmt"
  "golang.org/x/net/html"
  "net/http"
  "os"
)

const base_url = "https://www.pro-football-reference.com"

// Extract all hyperlinks from a webpage
func get_page_hrefs(page_url string) (page_hrefs []string) {
  resp, err := http.Get(page_url)
  if err != nil {
    fmt.Println("ERROR: Failed to crawl \"" + page_url + "\"")
    return
  }
  defer resp.Body.Close()

  page_tokens := html.NewTokenizer(resp.Body)
  for {
    tt := page_tokens.Next()
    switch {
    case tt == html.ErrorToken:
      // End of the document, we're done
        return
    case tt == html.StartTagToken:
      t := page_tokens.Token()

      isAnchor := t.Data == "a"
      if isAnchor {
        for _, a := range t.Attr {
          if a.Key == "href" {
            page_hrefs = append(page_hrefs, a.Val)
            break
          }
        }
      }
    }
  }

  return
}

func main() {
  team := os.Args[1]
  year := os.Args[2]

  team_url := fmt.Sprintf("%s/teams/%s/%s.htm", base_url, team, year)
  fmt.Println("Team URL: " + team_url)

  hrefs := get_page_hrefs(team_url)
  for _, h := range hrefs {
    fmt.Println(h)
  }
}
