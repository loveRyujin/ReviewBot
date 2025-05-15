package prompt

const DefaultLanguage = "English"

var languageMaps = map[string]string{
	"en":    DefaultLanguage,
	"zh-tw": "Traditional Chinese",
	"zh-cn": "Simplified Chinese",
	"ja":    "Japanese",
}

// GetLanguage returns the language name for the given language key,
// or the default language if the code is not recognized.
func GetLanguage(langKey string) string {
	if language, ok := languageMaps[langKey]; ok {
		return language
	}
	return DefaultLanguage
}
