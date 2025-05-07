package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
	"test/end2end/page"
	"testing"
)

func TestRank(t *testing.T) {
	testFunc := func(t *testing.T, driver selenium.WebDriver) {
		indexPage := page.Index{}
		indexPage.Init(driver)

		err := indexPage.OpenPage("")
		assert.NoError(t, err, "Не удалось открыть главную страницу")

		expectedText := "123a"
		expectedRank := 0.4
		expectedSimilarity := 0
		err = indexPage.InputText(expectedText)
		assert.NoError(t, err, "Не удалось ввести текст")

		var currentURL string
		err = indexPage.WaitWithTimeoutAndInterval(currentURL)
		assert.NoError(t, err, "Не перейти на страницу")

		summaryPage := page.Summary{}
		summaryPage.Init(driver)

		actualText, err := summaryPage.GetResultText()
		assert.NoError(t, err)
		assert.Equal(t, expectedText, actualText, "Text не совпадает")

		actualRank, err := summaryPage.GetResultRank()
		assert.NoError(t, err)
		assert.Equal(t, expectedRank, actualRank, "Rank не совпадает")

		actualSimilarity, err := summaryPage.GetResultSimilarity()
		assert.NoError(t, err)
		assert.Equal(t, expectedSimilarity, actualSimilarity, "Similarity не совпадает")
	}

	runTestForBrowser(t, "chrome", testFunc)
	runTestForBrowser(t, "firefox", testFunc)
}
