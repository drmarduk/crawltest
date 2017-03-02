package main

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/anaskhan96/soup"
)

func getNextLink(url string) (string, *soup.Node, error) {
	resp, err := downloadPage(url)
	if err != nil {
		return "", nil, err
	}
	doc := soup.HTMLParse(resp)

	// Element tag,(attribute key-value pair)
	e := doc.Find("span", "class", "next-button")
	if e == nil {
		return "", nil, errors.New("EOF")
	}
	return e.Find("a").Attrs()["href"], &doc, nil
}

func downloadPage(url string) (string, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	src, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

func extractImages(doc soup.Node) ([]string, error) {
	var result []string
	// <p class="title"><a class="title may-blank loggedin " data-event-action="title" href="http://imgur.com/bmdZt4F" tabindex="1" rel="" >What a cutie</a
	e := doc.FindAll("p", "class", "title")
	if e == nil {
		return result, errors.New("empty imageset")
	}
	for _, i := range e {
		img := i.Find("a").Attrs()["href"]
		if img == "" {
			continue
		}
		result = append(result, img)
	}

	return result, nil
}
