package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var url_prefix string = "http://"

func main() {
	for _, url := range os.Args[1:] {

		// Check url for http:// leader
		if !strings.HasPrefix(url, url_prefix) {
			url = url_prefix + url
		}

		// Get the url response as an io.reader
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error with %s: %s", url, err)
			os.Exit(0)
		}

		// Convert the reader to a text string
		body_text, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error with body %s: %s", url, err)
			os.Exit(0)
		}

		// Print the body to stdout
		fmt.Printf("Body of %s: %s", url, body_text)

	}
}
