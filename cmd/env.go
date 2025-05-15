package cmd

import (
	"github.com/spf13/viper"
)

var keyToEnv = map[string]string{
	"git.diff_file":        "GIT_DIFF_FILE",
	"git.max_input_size":   "GIT_MAX_INPUT_SIZE",
	"git.diff_unified":     "GIT_DIFF_UNIFIED",
	"git.exclude_list":     "GIT_EXCLUDE_LIST",
	"git.lang":             "GIT_LANG",
	"git.template_file":    "GIT_TEMPLATE_FILE",
	"git.template_string":  "GIT_TEMPLATE_STRING",
	"ai.socks":             "AI_SOCKS",
	"ai.api_key":           "AI_API_KEY",
	"ai.model":             "AI_MODEL",
	"ai.proxy":             "AI_PROXY",
	"ai.base_url":          "AI_BASE_URL",
	"ai.timeout":           "AI_TIMEOUT",
	"ai.max_tokens":        "AI_MAX_TOKENS",
	"ai.temperature":       "AI_TEMPERATURE",
	"ai.provider":          "AI_PROVIDER",
	"ai.skip_verify":       "AI_SKIP_VERIFY",
	"ai.headers":           "AI_HEADERS",
	"ai.top_p":             "AI_TOP_P",
	"ai.frequency_penalty": "AI_FREQUENCY_PENALTY",
	"ai.presence_penalty":  "AI_PRESENCE_PENALTY",
	"prompt.folder":        "PROMPT_FOLDER",
}

func init() {
	for key, env := range keyToEnv {
		viper.BindEnv(key, defaultEnvPrefix+"_"+env)
	}
}
