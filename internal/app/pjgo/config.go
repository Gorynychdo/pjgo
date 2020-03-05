package pjgo

type Config struct {
	Id       string `toml:"id"`
	Uri      string `toml:"uri"`
	Login    string `toml:"login"`
	Password string `toml:"password"`
	LogLevel uint   `toml:"log_level"`

	MakeCall  bool   `toml:"make_call"`
	RemoteUri string `toml:"remote_uri"`
}

func NewConfig() *Config {
	return &Config{
		Id:       "sip:test1@pjsip.org",
		Uri:      "sip:sip.pjsip.org",
		Login:    "test1",
		Password: "test1",
		LogLevel: 4,

		MakeCall:  true,
		RemoteUri: "sip:test1@pjsip.org",
	}
}
