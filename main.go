package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseURL = "https://discord.com/api/webhooks"

var lastRequest time.Time
var rateLimitResponse = []string{
	"Whoa there, speed racer! Let's take a breather and try again in a bit.",
	"Oops, looks like you've hit the fun limit! Come back later for more.",
	"Too much of a good thing, eh? Let's pace ourselves and try again soon.",
	"Looks like you've reached the end of the line! Please insert more tokens to continue.",
	"Hold your horses! You've exceeded the maximum number of requests. Try again later, cowboy!",
	"Slow down, tiger! You've hit the rate limit. Take a catnap and come back refreshed.",
	"Congratulations, you've achieved the rate limit! Your prize is waiting... just not right now.",
	"Uh-oh, you've tripped the alarm! The rate limit police are on their way. Lay low for a while.",
	"Looks like you've had too much fun for now. Come back when the coast is clear!",
	"You've reached the end of the all-you-can-request buffet. Time to let your requests digest!",
	"Ooh, sorry! You've hit the rate limit wall. Don't worry, it's not permanent. Try again later!",
	"Too many requests, too little time! Your request quota will refresh soon. Hang in there!",
	"Whoa, easy there! You've reached the request limit. Time to take a break and smell the roses.",
	"Well, well, well... if it isn't the rate limit! Your persistence is admirable, but let's take a breather.",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <port>")
		return
	}
	port := os.Args[1]
	http.HandleFunc("/", sendToWebhook)
	fmt.Println("Listening on localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func sendToWebhook(w http.ResponseWriter, r *http.Request) {
	rl := r.URL.Query().Get("rl")
	rateLimit := 0 * time.Second
	if rl != "" {
		rateLimitDuration, err := time.ParseDuration(rl)
		if err != nil {
			http.Error(w, "Invalid rate limit duration", http.StatusBadRequest)
			return
		}
		rateLimit = rateLimitDuration
	}
	if time.Since(lastRequest) < rateLimit {
		randomString := rateLimitResponse[rand.Intn(len(rateLimitResponse))]
		http.Error(w, randomString, http.StatusTooManyRequests)
		return
	}
	lastRequest = time.Now()
	path := r.URL.Path
	text := r.URL.Query().Get("text")
	ping := r.URL.Query().Get("ping")

	var pingText string
	if ping != "" {
		pingText = "<@&" + ping + ">\n"
	}

	if path == "/" || text == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	webhookPath := strings.TrimPrefix(path, "/")
	webhookURL := baseURL + "/" + webhookPath

	data := map[string]string{
		"content": pingText + text,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to send webhook", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Fprintf(w, "Message sent successfully")
	} else {
		http.Error(w, "Failed to send webhook", resp.StatusCode)
	}
}
