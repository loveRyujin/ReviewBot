package options

type Option interface {
	Apply(*Config)
}

type optionFunc func(*Config)

var _ Option = (*optionFunc)(nil)

func (f optionFunc) Apply(cfg *Config) {
	f(cfg)
}

func WithDiffUnified(diffUnified int) Option {
	return optionFunc(func(cfg *Config) {
		cfg.DiffUnified = diffUnified
	})
}

func WithExcludedList(excludedList []string) Option {
	return optionFunc(func(cfg *Config) {
		if len(excludedList) == 0 {
			return
		}
		cfg.ExcludedList = excludedList
	})
}

func WithIsAmend(isAmend bool) Option {
	return optionFunc(func(cfg *Config) {
		cfg.IsAmend = isAmend
	})
}

type Config struct {
	DiffUnified  int
	ExcludedList []string
	IsAmend      bool
}
