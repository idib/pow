package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	port := flag.Int("port", 3333, "Port number to run the server on")
	flag.Parse()

	if *port == 0 {
		fmt.Println("Usage: go run main.go -port <port>")
		return
	}

	if *port < 1 || *port > 65535 {
		fmt.Println("Invalid port number:", *port)
		return
	}

	//Of course, here you can make it so that it is read from a file or from any source.
	//I have provided an option, but did not implement it.
	quotes := []string{
		"You can never be overdressed or overeducated.",
		"Accept who you are. Unless you’re a serial killer.",
		"The best dreams happen when you’re awake.",
		"Success is not the key to happiness. Happiness is the key to success.",
		"Success is the child of audacity.",
	}

	stdRandomizer := rand.New(rand.NewSource(time.Now().UnixMicro()))

	//It is possible to use different provider implementations.
	randomQuoteProvider := NewRandomQuoteStatic(stdRandomizer, quotes)

	powGenerator := NewFindHashPOW(stdRandomizer, 6)

	server := NewServer(randomQuoteProvider, *port, powGenerator)
	server.run()
}
