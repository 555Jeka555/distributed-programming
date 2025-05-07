package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
	"lab9/page"
	"testing"
	"time"
)

func TestRank(t *testing.T) {
	testFunc := func(t *testing.T, driver selenium.WebDriver) {
		indexPage := page.Index{}
		indexPage.Init(driver)

		err := indexPage.OpenPage("/")
		assert.NoError(t, err, "Не удалось открыть главную страницу")

		expectedText := "123a"
		expectedRank := 0.75
		expectedSimilarity := 1
		err = indexPage.InputText(expectedText)
		assert.NoError(t, err, "Не удалось ввести текст")

		err = indexPage.SelectRegion("AE")
		assert.NoError(t, err, "Не удалось выбрать регион")

		err = indexPage.SubmitText()
		assert.NoError(t, err, "Не удалось отправить текст")

		var currentURL string
		currentURL, err = indexPage.WaitWithTimeoutAndInterval(currentURL)
		assert.NoError(t, err, "Не перейти на страницу")

		time.Sleep(3 * time.Second)
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
