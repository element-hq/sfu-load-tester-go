package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

type chromium struct {
	pw      *playwright.Playwright
	options playwright.BrowserTypeLaunchPersistentContextOptions
}

func newChromium(headless bool) (*chromium, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not run playwright: %v", err)
	}

	options := playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(headless),
		Args: []string{
			"--no-sandbox",
			"--use-fake-ui-for-media-stream",
			"--use-fake-device-for-media-stream",
			"--disable-web-security",
			"--allow-running-insecure-content",
			"--unsafely-treat-insecure-origin-as-secure",
			"--ignore-certificate-errors",
			"--autoplay-policy=no-user-gesture-required",
		},
	}

	return &chromium{pw, options}, nil
}

// Spawns the bots. Uses the array of supplied bots as the bots to run.
func (c *chromium) spawnBots(callURL string, bots []string, ctx context.Context) {
	for _, bot := range bots {
		fmt.Println("Spawning bot:", bot)
		_, err := c.spawnBot(callURL, bot)

		// Retry until success.
		for err != nil {
			fmt.Println("Bot failed to spawn:", bot, err)

			// Wait 3 seconds before trying again.
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
			}

			_, err = c.spawnBot(callURL, bot)
		}

		fmt.Println("Bot spawned:", bot)
	}
}

func (c *chromium) spawnBot(callURL string, bot string) (playwright.Page, error) {
	// Check if the persistent context folder (cache) exists already.
	// If it does, we can skip the log-in procedure. This checks if the `bot` folder
	// exists.
	skipLogin := false
	if _, err := os.Stat(bot); err == nil {
		// The folder exists, so we can skip the log-in procedure.
		skipLogin = true
	}

	// We use a persistent context, so that only the first launch of the tool is long.
	// Subsequent launches will be fast. We "cache" the context.
	context, err := c.pw.Chromium.LaunchPersistentContext(bot, c.options)
	if err != nil {
		return nil, fmt.Errorf("could not create context: %v", err)
	}

	// Fail relatively early (the default value is 30s).
	context.SetDefaultTimeout(5000)

	// Convenient function to handle Go boilerplate.
	returnError := func(err error, text string) (playwright.Page, error) {
		defer context.Close()
		return nil, fmt.Errorf("%s: %v", text, err)
	}

	page, err := context.NewPage()
	if err != nil {
		return returnError(err, "could not create page")
	}

	if _, err := page.Goto(callURL); err != nil {
		return returnError(err, "could not go to call URL")
	}

	if skipLogin {
		if err := page.Click("text=Join call now"); err != nil {
			return returnError(err, "could not click join call button")
		}

		return page, nil
	}

	// Otherwise, we need to log-in and then join the call.
	if err := page.Click("text='Log in'"); err != nil {
		return returnError(err, "could not click login button")
	}

	if err := page.Fill("[placeholder='Username']", bot); err != nil {
		return returnError(err, "could not fill username")
	}

	if err := page.Fill("input[placeholder=Password]", bot); err != nil {
		return returnError(err, "could not fill password")
	}

	if err := page.Click("text=Login"); err != nil {
		return returnError(err, "could not click login button")
	}

	if err := page.Click("text='Join call now'"); err != nil {
		return returnError(err, "could not click join call button")
	}

	return page, nil
}

func (c *chromium) close() error {
	if err := c.pw.Stop(); err != nil {
		return fmt.Errorf("could not close browser: %v", err)
	}

	return nil
}
