package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type chromium struct {
	browser playwright.Browser
}

func newChromium() (*chromium, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not run playwright: %v", err)
	}

	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
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

	browser, err := pw.Chromium.Launch(launchOptions)
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}

	return &chromium{browser}, nil
}

// Runs the bots until the stop signal is received.
// Uses the array of supplied bots as the bots to run.
func (c *chromium) runBots(callURL string, bots []string, ctx context.Context) error {
	go func() {
		<-ctx.Done()
		for _, browserCtx := range c.browser.Contexts() {
			browserCtx.Close()
		}
		c.browser.Close()
	}()

	for _, bot := range bots {
		page, err := c.spawnBot(callURL, bot)

		if err != nil {
			if !strings.Contains(err.Error(), "Timeout") {
				return err
			}

			// It's likely a rate-limitting issue, so we'll just wait a bit and try again.
			fmt.Printf("Hitting rate-limit, waiting to bypass the rate-limitter...:\n")

			select {
			case <-ctx.Done():
			case <-time.After(3 * time.Second):
			}

			// Try again, if we failed again, then it's not rate-limitting, but something else.
			page, err = c.spawnBot(callURL, bot)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Bot %s spawned\n", bot)
		defer page.Close()
	}

	// Wait for the stop signal.
	<-ctx.Done()

	return nil
}

func (c *chromium) spawnBot(callURL string, bot string) (playwright.Page, error) {
	context, err := c.browser.NewContext()
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
	if err := c.browser.Close(); err != nil {
		return fmt.Errorf("could not close browser: %v", err)
	}

	return nil
}
