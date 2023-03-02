package main

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type callConfig struct {
	callURL string
	room    string
}

type chromium struct {
	browser playwright.Browser
}

func newChromium() (*chromium, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not run playwright: %v", err)
	}

	options := playwright.BrowserTypeLaunchOptions{
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

	browser, err := pw.Chromium.Launch(options)
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}

	return &chromium{browser}, nil
}

func (c *chromium) run(config callConfig, numberBots int, stopSignal <-chan struct{}) error {
	pages := make([]playwright.Page, numberBots)
	for i := 0; i < numberBots; i++ {
		page, err := c.spawnBot(config, i)
		if err != nil {
			return fmt.Errorf("could not spawn bot: %v", err)
		}

		pages[i] = page
	}

	// Wait for the stop signal.
	<-stopSignal
	for _, page := range pages {
		if err := page.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chromium) spawnBot(config callConfig, index int) (playwright.Page, error) {
	context, err := c.browser.NewContext()
	if err != nil {
		return nil, fmt.Errorf("could not create context: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}

	if _, err := page.Goto(config.callURL); err != nil {
		return nil, fmt.Errorf("could not goto page: %v", err)
	}

	if err := page.Type("input#callName", config.room); err != nil {
		return nil, fmt.Errorf("could not type room: %v", err)
	}

	if err := page.Type("input#displayName", fmt.Sprintf("bot_%d", index)); err != nil {
		return nil, fmt.Errorf("could not type user: %v", err)
	}

	if err := page.Press("input#displayName", "Enter"); err != nil {
		return nil, fmt.Errorf("could not press enter: %v", err)
	}

	if _, err := page.WaitForNavigation(); err != nil {
		return nil, fmt.Errorf("could not wait for navigation: %v", err)
	}

	// Find the button that has a text of "Join call now" and press it.
	if err := page.Click("text=Join call now"); err != nil {
		return nil, fmt.Errorf("could not click join call now: %v", err)
	}

	return page, nil
}

func (c *chromium) close() error {
	if err := c.browser.Close(); err != nil {
		return fmt.Errorf("could not close browser: %v", err)
	}

	return nil
}
