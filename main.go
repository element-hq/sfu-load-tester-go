package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TODO: Read from cmd parameters.
var (
	config = callConfig{
		callURL: "https://pr805--element-call.netlify.app/",
		room:    "dcall",
	}
	numberBots = 10
)

func main() {
	// Creates a new chromium instance.
	botBrowser, err := newChromium()
	if err != nil {
		log.Fatalf("could not create chrome bot: %v", err)
	}

	// Close the browser when the app is done.
	defer botBrowser.close()

	// Handle SIGTERM in terminal, so that once we're done, we stop the bots.
	signalStop := make(chan struct{})
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		close(signalStop)
	}()

	// Run the bots until the stop signal is received.
	if err := botBrowser.run(config, numberBots, signalStop); err != nil {
		log.Fatalf("could not spawn bots: %v", err)
	}
}
