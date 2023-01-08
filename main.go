package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Message Channel Declaration
var ch = make(chan string)

func closeChannel() {
	close(ch)
	ch = nil
	fmt.Println("Connection Closed")
	fmt.Println("Channel Closed")
}
func main() {
	router := http.NewServeMux()
	router.HandleFunc("/event", sseHandler)
	router.HandleFunc("/time", timeHandler)
	log.Fatal(http.ListenAndServe(":8080", router))

}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer closeChannel() // Close Channel on exit

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client Disconnected")
			return
		}
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if ch != nil {
		ch <- time.Now().Format(time.RFC3339)
		fmt.Println("Message Sent")
	} else {
		fmt.Println("Channel Closed")
	}
}
