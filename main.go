package main

import (
	"log"
	"time"
	//"runtime"

	"./stat"
	"./api"
)

func main() {
	//runtime.GOMAXPROCS(2)

	go api.Start()

	go runThread()

	quit := make(chan bool)
	<-quit

}

func runThread() {
	for {
		/*
		log.Println("starting system screen update.")
		systemScreens := stat.UpdateSystemScreen()

		log.Println("reading active screen config.")
		go readActiveScreensConfig(&activeScreens, "active_screen.json")

		log.Println("checking screens.")
		go CheckScreens(activeScreens, systemScreens)
		*/

		_, err := stat.CollectSystemStats()
		if err != nil {
			log.Println(err)
		}

		_, err = stat.CollectStrSystemStats()
		if err != nil {
			log.Println(err)
		}

		/*
		strSystemStats, err := stat.CollectStrSystemStats()
		if err != nil {
			log.Println(err)
		}
		log.Println(strSystemStats)
		*/
		_, err = stat.CollectScreenStats()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(10 * time.Second)
	}
}
