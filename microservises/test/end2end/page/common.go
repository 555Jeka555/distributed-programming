package page

import (
	"fmt"
	"test/end2end/config"
	"time"

	"github.com/tebeka/selenium"
)

type Common struct {
	driver  selenium.WebDriver
	baseURL string
}

func (c *Common) Init(driver selenium.WebDriver) {
	c.driver = driver
	c.baseURL = config.BaseUrl
}

func (c *Common) OpenPage(url string) error {
	return c.driver.Get(c.baseURL + url)
}

func (c *Common) WaitWithTimeoutAndInterval(currentURL string) (string, error) {
	err := c.driver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		url, err := wd.CurrentURL()
		if err != nil {
			return false, err
		}
		currentURL = url
		return url != "", nil
	}, 10*time.Second, 500*time.Millisecond)

	return currentURL, err
}

func (c *Common) findElement(by, value string) (selenium.WebElement, error) {
	return c.waitForElement(by, value, 10*time.Second)
}

func (c *Common) waitForElement(by, value string, timeout time.Duration) (selenium.WebElement, error) {
	var elem selenium.WebElement
	var err error

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		elem, err = c.driver.FindElement(by, value)
		if err == nil {
			return elem, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("element not found after %v: %v", timeout, err)
}
