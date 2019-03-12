package main

import (
	"log"
	"time"
	"flag"
	"os"

	"./stat"
	"./api"
)

var options struct {
	config string
}

func main() {
	flag.StringVar(&options.config, "config", "", "name of config file json format")
	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal("No args set. At least set --config <config-filename>")
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
		time.Sleep(10 * time.Second)
	}
}
