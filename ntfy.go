package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Action struct {
	Action string `json:"action"`
	Label  string `json:"label"`
	Url    string `json:"url"`
}

type Alert struct {
	Topic    string    `json:"topic"`
	Message  string    `json:"message"`
	Markdown bool      `json:"markdown"`
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	Priority int       `json:"priority"`
	Attach   string    `json:"attach"`
	Filename string    `json:"filename"`
	Click    string    `json:"click"`
	Actions  []*Action `json:"actions"`
}

func SendNotification(alert *Alert) {
	alertBody, _ := json.Marshal(alert)
	req, err := http.NewRequest("POST", "https://ntfy.sh", bytes.NewBuffer(alertBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
