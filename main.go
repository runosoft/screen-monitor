package main

import (
	"log"
	"time"
	"flag"
	"os"

	"github.com/arsmine/screen-monitor/stat"
	"github.com/arsmine/screen-monitor/api"
)

var options struct {
	config string
	interval string
}

func main() {
	flag.StringVar(&options.config, "config", "", "name of config file json format")
	flag.StringVar(&options.interval, "interval", "", "check interval of general operations")
	flag.Parse()

	if len(os.Args) < 3 {
		log.Fatal("Not enough args set. set --config <config-filename> and --interval <duration>")
	}

	go api.Start()

	go runThread(options.config)

	quit := make(chan bool)
	<-quit
}

func runThread(config string) {
	for {
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

		interval, err := time.ParseDuration(options.interval)
		if err != nil {
			log.Fatalf("Couldn't parse interval")
		}
		time.Sleep(interval)
	}
}
