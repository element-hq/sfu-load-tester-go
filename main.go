package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

// TODO: Parse these parameters.
var (
	botUsers = []string{
		"gobotuser1",
		"gobotuser2",
		"gobotuser3",
		"gobotuser4",
		"gobotuser5",
		"gobotuser6",
		"gobotuser7",
		"gobotuser8",
		"gobotuser9",
		"gobotuser10",
		"gobotuser11",
		"gobotuser12",
		"gobotuser13",
		"gobotuser14",
		"gobotuser15",
		"gobotuser16",
		"gobotuser18",
		"gobotuser19",
		"gobotuser20",
	}
	callURL = "https://pr805--element-call.netlify.app/room/#dcall1:call.ems.host"
)

func main() {
	// Creates a new chromium instance.
	botBrowser, err := newChromium()
	if err != nil {
		log.Fatalf("could not create chrome bot: %v", err)
	}

	// Close the browser when the app is done.
	defer botBrowser.close()

	// Creates a context that will be used to cancel the bots.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Runs the bots until the context is closed.
	if err := botBrowser.runBots(callURL, botUsers, ctx); err != nil {
		log.Fatalf("could not spawn bots: %v", err)
	}
}
