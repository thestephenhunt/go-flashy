package main

import (
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
	NewClients     chan chan string
	ClosingClients chan chan string
	Clients        map[chan string]bool
	Messages       chan string
}

func (broker *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)

	broker.NewClients <- messageChan

	notify := r.Context().Done()

	go func() {
		<-notify
		broker.ClosingClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		msg, open := <-messageChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}

	log.Println("Finished HTTP request at ", r.URL.Path)
}

func (broker *Broker) Start() {
	go func() {
		for {
			select {
			case s := <-broker.NewClients:
				broker.Clients[s] = true
				log.Printf("Client added. %d registered clients", len(broker.Clients))
			case s := <-broker.ClosingClients:
				delete(broker.Clients, s)
				close(s)
				log.Printf("Removed client. %d registered clients", len(broker.Clients))
			case msg := <-broker.Messages:
				for s := range broker.Clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(broker.Clients))
			}
		}
	}()
}
