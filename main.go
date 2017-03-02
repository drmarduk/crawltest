package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {

	Downloader(4)
	reddit := "https://www.reddit.com/r/ScarlettJohansson/new/"

	log.Println("Get: ", reddit)
	next, src, err := getNextLink(reddit)
	if err != nil {
		panic(err)
	}
	// do initial picture extraction on src
	for {
		log.Println("Get: ", next)
		next, src, err = getNextLink(next)
		if err != nil || src == nil {
			log.Println("end reached")
			break
		}

		images, err := extractImages(*src)
		if err != nil {
			log.Printf("Error while extracting images: %s\n", err)
			continue
		}

		//time.Sleep(1 * time.Second)
		for _, i := range images {
			WorkQueue <- i
		}

	}

	wg.Wait()
}

// WorkQueue holds all incoming links
var WorkQueue = make(chan string, 100)

// Downloader asdf
func Downloader(n int) {
	go func() {
		for {
			select {
			case url := <-WorkQueue:
				go downloadToDisk(url)
			}
		}
	}()
}

// download stuff
func fileexists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

// download a file to disk, does overwrite existing files
func downloadToDisk(url string) {
	wg.Add(1)
	defer wg.Done()

	x := strings.Split(url, "/")
	fn := x[len(x)-1] // get filename

	if fileexists(fn) {
		return
	}

	out, err := os.Create(fn)
	if err != nil {
		return
	}
	defer out.Close()
	download(out, url)
	log.Println("Downloaded: ", url)
}

func download(dst io.Writer, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
