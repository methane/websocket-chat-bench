package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"sync"
)

var (
	origin     string
	url        string
	numClients int
	numMessage int
)

func init() {
	flag.StringVar(&origin, "origin", "http://localhost", "Origin header")
	flag.StringVar(&url, "url", "ws://localhost:5000/chat", "Target URL.")
	flag.IntVar(&numClients, "n", 1, "Number of clients.")
	flag.IntVar(&numMessage, "m", 100, "Number of messages per client.")
}

func client(wg *sync.WaitGroup, ws *websocket.Conn) {
	defer wg.Done()
	go func() {
		defer ws.Close()
		for i := 0; i < numMessage; i++ {
			err := websocket.Message.Send(ws, "Hello")
			if err != nil {
				log.Println(err)
				break
			}
		}
	}()

	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Received %#v\n", msg)
	}
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go client(&wg, ws)
	}
	wg.Wait()
}
