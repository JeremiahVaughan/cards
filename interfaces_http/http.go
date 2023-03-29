package interfaces_http

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type logWriter struct {
}

func Http() {
	resp, err := http.Get("https://google.com")
	if err != nil {
		log.Fatalf("error, when making a call to google.com: %v", err)
	}

	//io.Copy(os.Stdout, resp.Body)
	//if err != nil {
	//	log.Fatalf("error, when attempting to read data from response body: %v", err)
	//}
	//log.Printf(string(data))

	lw := logWriter{}
	io.Copy(lw, resp.Body)
	if err != nil {
		log.Fatalf("error, when attempting to read data from response body: %v", err)
	}
}

func (lw logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("just wrote this many bytes:", len(bs))
	return len(bs), nil
}
