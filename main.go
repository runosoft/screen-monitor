package main

import(
	"os/exec"
	"fmt"
	"strings"
	"log"
	"encoding/json"
	"os"
	"path/filepath"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLStruct struct {
	Storage StorageConfig `json:"db"`
}

type StorageConfig struct {
	Driver string `json:"driver"`
	Name string `json:"name"`
}

/* Holds screen informations from DB */
type DBScreens struct {
	PID string
	Name string
}

/* Holds screen informations from command */
type SystemScreens struct {
	PID string
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
		log.Fatalf("Error when opening %s: %s\n", configFilename ,err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fataf("Config error for %s: %s\n", configFilename, err.Error())
	}
}

func readActiveScreensConfig(activeScreen *ActiveScreens, configFilename string) {
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Reading Active Screens: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalf("Error when opening %s: %s\n", configFilename ,err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&activeScreen); err != nil {
		log.Fataf("Config error for %s: %s\n", configFilename, err.Error())
	}
}

func main() {
	readConfig(&config, "config.json")
	readActiveScreensConfig(&activeScreens, "active_screen.json")

	db, err := sql.Open(config.Storage.Driver, config.Storage.Name)
	if err != nil {
		log.Fatal(err)
	}

	screenStruct := updateScreenList(db)
	log.Println(screenStruct)

	checkScreens()
}

func checkScreens(screens []DBScreens) {
	for i := 0; i < len(screens); i++ {
		
	}
}

/* Parses output of screen -ls command to SystemScreens 
 * struct. */
func updateSystemScreen() []SystemScreens{
	out, err := exec.Command("screen", "-ls").Output()
	if err != nil {
		log.Fatalf("Error while executing screen -ls command: %s\n", err)
	}
	screenOut := string(out)
	screenOutLineArr := strings.Split(screenOut, "\n")

	var screenStructArr []SystemScreens
	var pureScreen []string

	for i := 1; i < len(screenOutLineArr)-2; i++ {
		pureScreen = append(pureScreen, screenOutLineArr[i])
	}

	for i := 0; i < len(pureScreen); i++ {
		screenOut = strings.TrimSpace(pureScreen[i])
		screenOutArr := strings.Split(screenOut, "\t")
		screenOut = screenOutArr[0]
		screenOutArr = strings.Split(screenOut, ".")
		screenPID := screenOutArr[0]
		screenName := screenOutArr[1]

		systemScreens := SystemScreens{
			PID: screenPID,
			Name: screenName,
		}

		screenStructArr = append(screenStructArr, SystemScreens)
	}
	return screenStructArr
}

/* Gets screen PIDs, names from DB and parse them to
 * the DBScreens struct. */
func updateDBScreen(db *sql.DB) []DBScreens {
	var screenStructArr []DBScreens

	rows, err := db.Query("SELECT PID, screen_name FROM ScreenInfo")
	if err != nil {
		log.Printf("Error while getting PID and screen_name from DB: %s\n", err)
	}

	var dbScreenPID string
	var dbScreenName string

	for rows.Next() {
		err = rows.Scan(&dbScreenPID, &dbScreenName)
		if err != nil {
			log.Printf("Error while scanning screen info from DB: %s\n", err)
		}

		dbScreens := DBScreens{
			PID: dbScreenPID,
			Name: dbScreenName,
		}
		screenStructArr = append(screenStructArr, screenStruct)
	}
	return screenStructArr
}

func RowExists(db *sql.DB, pid string) bool {
	existQuery := "SELECT exists(SELECT PID FROM ScreenInfo WHERE PID=?)"
	var exists bool

	err := db.QueryRow(existQuery, pid).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error when getting existence of pid from DB: %s\n", err)
		return false
	}
	return exists
}
