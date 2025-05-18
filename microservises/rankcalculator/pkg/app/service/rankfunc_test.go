package service

import (
	"testing"
)

func TestRankCalculator_Parametric(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected float64
	}

	tests := []testCase{
		{
			name:     "Только латиница",
			input:    "abcXYZ",
			expected: 0,
		},
		{
			name:     "Только кириллица",
			input:    "абвгдЕЁ",
			expected: 0,
		},
		{
			name:     "Только цифры",
			input:    "123456",
			expected: 1,
		},
		{
			name:     "Только спецсимволы",
			input:    "!@#$%^",
			expected: 1,
		},
		{
			name:     "Смешанная строка",
			input:    "abc123!@#",
			expected: 0.6666666666666667,
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: 0,
		},
		{
			name:     "Строка с пробелами",
			input:    "a b c",
			expected: 0.4,
		},
		{
			name:     "Один Emoji",
			input:    "😀",
			expected: 1,
		},
		{
			name:     "Несколько Emoji и буквы",
			input:    "a😀b🐶c",
			expected: 0.4,
		},
		{
			name:     "Китайские иероглифы",
			input:    "你好abc",
			expected: 0.4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := calcRank(tc.input)
			if (actual-tc.expected) > 1e-9 || (tc.expected-actual) > 1e-9 {
				t.Errorf("Rank(%q) = %v, ожидалось %v", tc.input, actual, tc.expected) // TODO задействовать testifyassert
			}
		})
	}
}
