package main

// This program uses a RequestPool to log tweet IDs
import (
	"encoding/json"
	"fmt"
	pool "github.com/alexbotello/requests/pool"
	"io/ioutil"
	"log"
	"net/http"
)

type TwitterUser struct {
	username string
}

// to store the tweets from response
var tweetList []string

// The TwitterUser type can now be considered a Requestor as it
// implements a Request method.
func (t TwitterUser) Request() {
	twitterToken := "blahblahsometokenhere"
	authToken := "Bearer " + twitterToken

	// endpoint to get the 5 most recent tweets based on username
	url := "https://api.twitter.com/1.1/statuses/user_timeline.json?count=5&screen_name="
	apiRoute := fmt.Sprintf("%v%v", url, t.username)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiRoute, nil)
	req.Header.Add("Authorization", authToken)

	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var responseData []interface{}
	_ = json.Unmarshal(body, &responseData)
	for _, record := range responseData {
		if rec, ok := record.(map[string]interface{}); ok {
			id := rec["id_str"].(string)
			tweetList = append(tweetList, id)
		}
	}
}

// Lets make an array of Requestors based on different Twitter users we would
// like to get tweets for
func createRequestors() []pool.Requestor {
	usernames := []string{"BarackObama", "justinbieber", "katyperry", "rihanna", "ladygaga", "mastodonmusic"}

	var requestors []pool.Requestor
	for _, name := range usernames {
		t := &TwitterUser{username: name}
		requestors = append(requestors, t)
	}
	return requestors
}

func main() {
	log.Print("Starting request pool...")
	requestors := createRequestors()

	rp := pool.NewRequestPool(requestors)
	rp.Start() // this will block until all request workers have finished

	for _, tweet := range tweetList {
		log.Printf("Found Tweet: %v", tweet)
	}
	log.Print("Exiting...")
}
