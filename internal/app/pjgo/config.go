package pjgo

type Config struct {
	Id       string `toml:"id"`
	Uri      string `toml:"uri"`
	Login    string `toml:"login"`
	Password string `toml:"password"`
}

func NewConfig() *Config {
	return &Config{
		Id:       "sip:test1@pjsip.org",
		Uri:      "sip:sip.pjsip.org",
		Login:    "test1",
		Password: "test1",
	}
}
