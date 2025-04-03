package ai

type Provider string

const (
	OpenAI    Provider = "openai"
	Anthropic Provider = "anthropic"
	DeepSeek  Provider = "deepseek"
)

func (p Provider) String() string {
	return string(p)
}
