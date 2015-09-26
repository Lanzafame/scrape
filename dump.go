package main

import (
  "fmt"
  "net/http"
  "os"
)

func crawl(url string) {
  resp, err := http.Get(url)

  b := resp.Body
  return b
}

func main() {
  seedUrls := os.Args[1:]

  html := crawl(url)

  fmt.Println(html)

}

