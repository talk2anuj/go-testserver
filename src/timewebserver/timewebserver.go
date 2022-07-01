package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Welcome</h1>")
}

func getTimeFeed() string {
	fmt.Println("get Current time from Sever: ", timeProviderUrl)
	resp, err := http.Get(timeProviderUrl)
	if err != nil {
		// TODO: this path needs to be handled better
		fmt.Println("Got error from time provider server. Make sure it is running")
		return ""
	} else {
		fmt.Println("Got response from server")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// TODO : handle failure to decode response
			fmt.Println("Got error reading response")
			return ""
		} else {
			fmt.Printf("Got time: %s\n", body)
			newVal := fmt.Sprintf("%s", body)
			return newVal
		}
	}
}

func gettime(w http.ResponseWriter, r *http.Request) {
	// get time from local cache
	lastTime, found := timeCache.Get("lastTime")
	lastTimeAsString := lastTime.(string)
	if found {
		fmt.Println("Cache hit. Data: ", lastTimeAsString)
	} else {
		// not found in cache, get from server and update cache
		lastTimeAsString = getTimeFeed()
		timeCache.Set("lastTime", lastTimeAsString, cache.DefaultExpiration)
	}

	fmt.Fprintf(w, lastTimeAsString)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Health check</h1>")
}

// initalize cache
var timeCache = cache.New(5*time.Minute, 10*time.Minute)
var timeProviderUrl = "http://localhost:3000/gettime"

func main() {
	// create cron job to update cache
	cacheUpdateJob := cron.New()

	// run every 5 seconds
	cacheUpdateJob.AddFunc("*/10 * * * * *", func() {
		fmt.Println("Getting new data to update cache")
		newVal := getTimeFeed()
		fmt.Println("updating cache with: ", newVal)
		timeCache.Set("lastTime", newVal, cache.DefaultExpiration)
	})

	cacheUpdateJob.Start()

	// start server
	http.HandleFunc("/", index)
	http.HandleFunc("/health_check", healthCheck)
	http.HandleFunc("/gettime", gettime)
	fmt.Println("Server starting...")
	http.ListenAndServe(":3001", nil)

	fmt.Println("Stopping Server")
	cacheUpdateJob.Stop()

}
