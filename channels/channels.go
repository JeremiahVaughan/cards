package channels

import (
	"fmt"
	"net/http"
	"time"
)

func Channels() {
	sitesToCheck := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	// option 1
	//wg := sync.WaitGroup{}
	//wg.Add(len(sitesToCheck))
	//for _, s := range sitesToCheck {
	//	go checkLink(s)
	//}
	//wg.Wait()

	//	option 2
	c := make(chan string)
	for _, s := range sitesToCheck {
		go checkLink(s, c)
	}

	// Syntax 1
	//for {
	//	go checkLink(<-c, c)
	//}
	// Syntax 2
	for l := range c {
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Printf("error, when sending get request to %s. Error: %v\n", link, err)
		c <- link
		return
	}

	fmt.Printf("%s is up\n", link)
	c <- link
}
