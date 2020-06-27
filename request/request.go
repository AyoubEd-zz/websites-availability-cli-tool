package request

import (
	"fmt"
	"net/http"
	"time"
)

// Send performs a request to the given URL
func Send(t time.Time, url string) error {
	fmt.Println(url)
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	_ = resp
	fmt.Println("response from ->", t, " ", url)
	return nil
}
