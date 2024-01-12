package ucodesdk

type Config struct {
	AppId      string
	BaseURL    string
	TableSlug  string
	BotToken   string
	AccountIds []string
}

func (cfg *Config) SetAppId(appId string) {
	cfg.AppId = appId
}

func (cfg *Config) SetBaseUrl(url string) {
	cfg.BaseURL = url
}

func (cfg *Config) SetTableSlug(slug string) {
	cfg.TableSlug = slug
}

func (cfg *Config) SetBotToken(token string) {
	cfg.BotToken = token
}
