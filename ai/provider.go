package ai

type Provider string

const (
	OpenAI    Provider = "openai"
	Anthropic Provider = "anthropic"
	DeepSeek  Provider = "deepseek"
	Gemini    Provider = "gemini"
)

func (p Provider) String() string {
	return string(p)
}
