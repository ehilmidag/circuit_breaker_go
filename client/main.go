package main

import (
	"fmt"
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
	"time"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var circuitSettings gobreaker.Settings
	circuitSettings.Name = "cb example"
	circuitSettings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.4
	}
	circuitSettings.Timeout = 10 * 100000
	circuitSettings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		fmt.Println(name, " state change from ", from, " to ", to)
	}
	cb = gobreaker.NewCircuitBreaker(circuitSettings)
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("error come from here")
			fmt.Println("Error: ", err)
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}
		return body, nil
	})
	if err != nil {
		return nil, err
	}
	return body.([]byte), nil
}

func main() {
	corruptedUrl := "https://localhost:8081"
	correctUrl := "http://localhost:8080"

	var err error
	var body []byte

	for i := 0; i < 20; i++ {
		fmt.Println("Request: ", i)
		body, err = Get(corruptedUrl)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println("Body: ", string(body))
		if i > 14 {
			corruptedUrl = correctUrl
		}
		time.Sleep(3000)
	}
}
