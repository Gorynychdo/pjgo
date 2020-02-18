package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Gorynychdo/pjgo/internal/app/pjgo"
	"log"
	"os"
	"os/signal"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/pjgo.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := pjgo.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	sipUser := pjgo.NewSipUser()
	sipService := pjgo.NewSipService(sipUser)
	sipUser.SipService = sipService

	sipService.RegisterAccount(config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
