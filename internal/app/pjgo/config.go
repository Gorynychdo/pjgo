package pjgo

type Config struct {
	Id       string `toml:"id"`
	Uri      string `toml:"uri"`
	Login    string `toml:"login"`
	Password string `toml:"password"`
	LogLevel uint   `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		Id:       "sip:test1@pjsip.org",
		Uri:      "sip:sip.pjsip.org",
		Login:    "test1",
		Password: "test1",
		LogLevel: 4,
	}
}
