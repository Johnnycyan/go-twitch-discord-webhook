package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const baseURL = "https://discord.com/api/webhooks"

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
