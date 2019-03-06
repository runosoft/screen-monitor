package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"runtime"
)

type SQLStruct struct {
	Storage StorageConfig `json:"db"`
}

type StorageConfig struct {
	Driver string `json:"driver"`
	Name   string `json:"name"`
}

/* Holds screen informations from DB */
type DBScreens struct {
	PID  string
	Name string
}

/* Holds screen informations from command */
type SystemScreens struct {
	PID  string
	Name string
}

/* Holds name of the screen that
 * we want to check whether active or not */
type ActiveScreens struct {
	Names []string `json:"activeScreen"`
}

var config SQLStruct
var activeScreens ActiveScreens

func readConfig(cfg *SQLStruct, configFileName string) {
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalf("Error when opening %s: %s\n", configFile, err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fatalf("Config error for %s: %s\n", configFile, err.Error())
	}
}

func readActiveScreensConfig(activeScreen *ActiveScreens, configFileName string) {
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Reading Active Screens: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalf("Error when opening %s: %s\n", configFileName, err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&activeScreen); err != nil {
		log.Fatalf("Config error for %s: %s\n", configFile, err.Error())
	}
}

func main() {
	readConfig(&config, "config.json")

	runtime.GOMAXPROCS(2)

	go runThread()

	quit := make(chan bool)
	<-quit

}

func runThread() {
	for {
		log.Println("starting system screen update.")
		systemScreens := updateSystemScreen()

		log.Println("reading active screen config.")
		go readActiveScreensConfig(&activeScreens, "active_screen.json")

		log.Println("checking screens.")
		go checkScreens(activeScreens, systemScreens)
		time.Sleep(30 * time.Second)
	}
}

func checkScreens(activeScreensCfg ActiveScreens, systemScreens []string) {
	for _, value := range activeScreensCfg.Names {
		exists := contains(systemScreens, value)
		if exists {
			log.Printf("%s is running.\n", value)
		} else {
			sendCrashMessage(value)
			log.Printf("%s is not running\n", value)
		}
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func sendCrashMessage(screenName string) {
	log.Printf("%s is crashed.\n", screenName)
}

/* Parses output of screen -ls command to SystemScreens
 * struct. */
func updateSystemScreen() []string {
	out, err := exec.Command("screen", "-ls").Output()
	if err != nil {
		log.Fatalf("Error while executing screen -ls command: %s\n", err)
	}
	screenOut := string(out)
	screenOutLineArr := strings.Split(screenOut, "\n")

	var sysScreenNames []string
	var pureScreen []string

	for i := 1; i < len(screenOutLineArr)-2; i++ {
		pureScreen = append(pureScreen, screenOutLineArr[i])
	}

	for i := 0; i < len(pureScreen); i++ {
		screenOut = strings.TrimSpace(pureScreen[i])
		screenOutArr := strings.Split(screenOut, "\t")
		screenOut = screenOutArr[0]
		screenOutArr = strings.Split(screenOut, ".")
		screenName := screenOutArr[1]

		sysScreenNames = append(sysScreenNames, screenName)
	}
	return sysScreenNames
}
