package main

import (
	"fmt"
	"os"

	"github.com/pluming/aurora/config"
	"github.com/pluming/aurora/lib/logger"
	"github.com/pluming/aurora/server"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

const Name = "aurora"

var defaultProperties = &config.ServerProperties{
	Bind:           "0.0.0.0",
	Port:           6389,
	AppendOnly:     false,
	AppendFilename: "",
	MaxClients:     1000,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       Name,
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFilename := os.Getenv("CONFIG")
	if configFilename == "" {
		if fileExists("redis.conf") {
			config.SetupConfig("redis.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupConfig(configFilename)
	}

	err := server.ListenAndServeWithSignal(&server.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	})
	if err != nil {
		logger.Error(err)
	}
}
