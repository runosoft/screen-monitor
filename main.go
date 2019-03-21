package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arsmine/screen-monitor/api"
	"github.com/arsmine/screen-monitor/config"
	"github.com/arsmine/screen-monitor/stat"
)

var options struct {
	config   string
	interval string
}

var mainCfg config.MainConfig

func readConfig(cfg *config.MainConfig, configFileName string) {
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}
}

func main() {
	flag.StringVar(&options.config, "config", "", "name of config file json format")
	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal("Not enough args set. set --config <config-filename>.")
	}

	readConfig(&mainCfg, options.config)
	go api.Start(&mainCfg)

	go runThread(options.config)

	quit := make(chan bool)
	<-quit
}

func runThread(config string) {
	for {
		readConfig(&mainCfg, config)

		_, err := stat.CollectSystemStats()
		if err != nil {
			log.Println(err)
		}

		_, err = stat.CollectStrSystemStats()
		if err != nil {
			log.Println(err)
		}

		_, err = stat.CollectScreenStats(config)
		if err != nil {
			log.Println(err)
		}

		interval, err := time.ParseDuration(mainCfg.Interval)
		if err != nil {
			log.Fatalf("Couldn't parse interval")
		}

		time.Sleep(interval)
	}
}
