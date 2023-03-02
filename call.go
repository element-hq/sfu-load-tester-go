package main

import (
	"context"
	"fmt"
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
		Headless: playwright.Bool(true),
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
	pages := make([]playwright.Page, len(bots))

	go func() {
		<-ctx.Done()
		c.browser.Close()
	}()

	for i, bot := range bots {
		page, err := c.spawnBot(callURL, bot)
		if err != nil {
			return fmt.Errorf("could not spawn bot: %v", err)
		}

		pages[i] = page
		defer page.Close()

		// Sleep to bypass the rate-limiting.
		select {
		case <-time.After(3 * time.Second):
		case <-ctx.Done():
		}
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

	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}

	if _, err := page.Goto(callURL); err != nil {
		return nil, fmt.Errorf("could not goto page: %v", err)
	}

	if err := page.Click("text='Log in'"); err != nil {
		return nil, fmt.Errorf("could not click login: %v", err)
	}

	if err := page.Fill("[placeholder='Username']", bot); err != nil {
		return nil, fmt.Errorf("could not fill username: %v", err)
	}

	if err := page.Fill("input[placeholder=Password]", bot); err != nil {
		return nil, fmt.Errorf("could not fill password: %v", err)
	}

	if err := page.Click("text=Login"); err != nil {
		return nil, fmt.Errorf("could not click login: %v", err)
	}

	if err := page.Click("text='Join call now'"); err != nil {
		return nil, fmt.Errorf("could not press enter: %v", err)
	}

	return page, nil
}

func (c *chromium) close() error {
	if err := c.browser.Close(); err != nil {
		return fmt.Errorf("could not close browser: %v", err)
	}

	return nil
}
