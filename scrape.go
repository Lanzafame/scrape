package main

import (
  "fmt"
  "net/http"
  "golang.org/x/net/html"
  "os"
  "strings"
)

// helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
  // iterate over all of the Token's attributes until we find an "href"
  for _, a := range t.Attr {
	if a.Key == "href" {
	  href = a.Val
	  ok = true
	}
  }

  // "bare" return will return the variables (ok, href) as defined in
  // the function definition
  return
}

// extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
  resp, err := http.Get(url)

  defer func() {
	// notify that we're done after this function
	chFinished <- true
  }()

  if err != nil {
	fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
	return
  }

  b := resp.Body
  defer b.Close() // close Body when the function returns

  z := html.NewTokenizer(b)

  for {
	tt := z.Next()

	switch {
	case tt == html.ErrorToken:
	  // end of the document, we're done
	  return
	case tt == html.StartTagToken:
	  t := z.Token()

	  // check if the token is an <a> tag
	  isAnchor := t.Data == "a"
	  if !isAnchor {
		continue
	  }

	  // extract the href value, if there is one
	  ok, url := getHref(t)
	  if !ok {
		continue
	  }

	  // make sure the url begins in http**
	  hasProto := strings.Index(url, "https") == 0
	  if hasProto {
		ch <- url
	  }
	}
  }
}

func main() {
  foundUrls := make(map[string]bool)
  seedUrls := os.Args[1:]

  // channels
  chUrls := make(chan string)
  chFinished := make(chan bool)

  // kick off the crawl process (concurrently)
  for _, url := range seedUrls {
	go crawl(url, chUrls, chFinished)
  }

  // subscribe to both channels
  for c := 0; c < len(seedUrls); {
	select {
	case url := <-chUrls:
	  foundUrls[url] = true
	case <-chFinished:
	  c++
	}
  }

  // we're done! print the results...
  
  fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

  for url, _ := range foundUrls {
	fmt.Println(" - " + url)
  }

  close(chUrls)
}


