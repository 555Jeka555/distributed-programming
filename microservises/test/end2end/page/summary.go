package page

import (
	"strconv"

	"github.com/tebeka/selenium"
)

const (
	resultTextCSSSelector       = ".result-content .result-item:nth-child(1) .result-value"
	resultRankCSSSelector       = ".result-content .result-item:nth-child(2) .result-value"
	resultSimilarityCSSSelector = ".result-content .result-item:nth-child(3) .result-value"
)

type Summary struct {
	Common
}

func (s *Summary) GetResultText() (string, error) {
	return s.getResultValue(resultTextCSSSelector)
}

func (s *Summary) GetResultRank() (float64, error) {
	val, err := s.getResultValue(resultRankCSSSelector)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(val, 64)
}

func (s *Summary) GetResultSimilarity() (int, error) {
	val, err := s.getResultValue(resultSimilarityCSSSelector)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

func (s *Summary) getResultValue(selector string) (string, error) {
	elem, err := s.findElement(selenium.ByCSSSelector, selector)
	if err != nil {
		return "", err
	}

	text, err := elem.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}
