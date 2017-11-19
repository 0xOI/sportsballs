package main

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/yhat/scrape"
	"net/http"
	"os"
	"strings"
)

const base_url = "https://www.pro-football-reference.com"

func get_page_rootnode(page_url string) (root *html.Node) {
  // request and parse the front page
  resp, err := http.Get(page_url)
  if err != nil {
    panic(err)
  }
  root, err = html.Parse(resp.Body)
  if err != nil {
    panic(err)
  }

  return
}

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

// get team record for entire season
func get_team_record(team string, year string) (team_record string) {
  team_url := fmt.Sprintf("%s/teams/%s/%s.htm", base_url, team, year)

  root := get_page_rootnode(team_url)

  matcher := func(n *html.Node) bool {
    return scrape.Attr(n, "data-stat") == "team_record"
  }
  x := scrape.FindAllNested(root, matcher)
  for _, x := range x {
    if scrape.Text(x) != "" {
      team_record = scrape.Text(x)
    }
  }

  return
}

// get head coach of team
func get_team_coach(team string, year string) (team_coach string) {
  team_url := fmt.Sprintf("%s/teams/%s/%s.htm", base_url, team, year)

  root := get_page_rootnode(team_url)

  extractor := func(values []string) string {
    if strings.Contains(values[0], "Coach:") {
      return values[2]
    }

    return ""
  }

  x := scrape.FindAll(root, scrape.ByTag(atom.P))
  for _, x := range x {
    t := scrape.TextJoin(x, extractor)
    if t != "" {
      team_coach = t
    }
  }

  return
}

func main() {
  team := os.Args[1]
  year := os.Args[2]

  fmt.Println("Team Record: " + get_team_record(team, year))
  fmt.Println("Team Coach: " + get_team_coach(team, year))
}
