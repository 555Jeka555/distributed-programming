package page

import (
	"fmt"

	"github.com/tebeka/selenium"
)

const (
	textAreaCSSSelector     = "textarea[name='text']"
	selectRegionCSSSelector = "select[name='country']"
	optionRegionXPath       = ".//option[@value='%s']"
	submitTextCSSSelector   = "input[type='submit']"
)

type Index struct {
	Common
}

func (i *Index) InputText(text string) error {
	textArea, err := i.findElement(selenium.ByCSSSelector, textAreaCSSSelector)
	if err != nil {
		return err
	}

	if err := textArea.Clear(); err != nil {
		return err
	}

	return textArea.SendKeys(text)
}

func (i *Index) SelectRegion(region string) error {
	selectElement, err := i.findElement(selenium.ByCSSSelector, selectRegionCSSSelector)
	if err != nil {
		return err
	}

	option, err := selectElement.FindElement(selenium.ByXPATH, fmt.Sprintf(optionRegionXPath, region))
	if err != nil {
		return err
	}

	return option.Click()
}

func (i *Index) SubmitText() error {
	submitElement, err := i.findElement(selenium.ByCSSSelector, submitTextCSSSelector)
	if err != nil {
		return err
	}

	return submitElement.Click()
}
