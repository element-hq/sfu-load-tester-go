package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

// TODO: Make them as parameters. Or register them dynamically.
var preRegisteredBots = []string{
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
	"gobotuser17",
	"gobotuser18",
	"gobotuser19",
	"gobotuser20",
}

func main() {
	callURL := flag.String("call-url", "", "The full URL to the Element Call.")
	numBots := flag.Int("num-bots", 0, "The number of bots to spawn.")
	headless := flag.Bool("headless", false, "Whether to run the browser in headless mode.")

	flag.Parse()

	if *callURL == "" || *numBots == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *numBots > len(preRegisteredBots) {
		log.Fatal("The number of bots is greater than the number of pre-registered bots.")
	}

	// Creates a new chromium instance.
	botBrowser, err := newChromium(*headless)
	if err != nil {
		log.Fatalf("could not create chrome bot: %v", err)
	}

	// Close the browser when the app is done.
	defer botBrowser.close()

	// Creates a context that will be used to cancel the bots.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Spawns the bots. The bots will run until `ctx` is cancelled or the `close()` is called.
	botBrowser.spawnBots(*callURL, preRegisteredBots[:*numBots], ctx)

	<-ctx.Done()
	fmt.Println("Shutting down bots...")
}
