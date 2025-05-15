package prompt

import "testing"

func Test_GetLanguage(t *testing.T) {
	testCases := []struct {
		langKey  string
		expected string
	}{
		{"en", "English"},
		{"zh-tw", "Traditional Chinese"},
		{"zh-cn", "Simplified Chinese"},
		{"ja", "Japanese"},
		{"invalid-key", "English"}, // Unknown language, should default to English
	}

	for _, tc := range testCases {
		t.Run(tc.langKey, func(t *testing.T) {
			result := GetLanguage(tc.langKey)
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}
