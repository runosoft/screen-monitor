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
	"strconv"

	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/mackerelio/go-osstat/disk"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/loadavg"
)

type PureSystemStat struct {
	MemoryTotal uint64
	MemoryFree uint64
	MemoryAvailable uint64
	SwapTotal uint64
	SwapFree uint64
	CPUSystem uint64
	CPUIdle uint64
	UpTime time.Duration
	DiskStats []disk.Stats
	NetworkStats []network.Stats
	LoadAvg1 float64
	LoadAvg5 float64
	LoadAvg15 float64
}

type StringSystemStat struct {
	StrMemoryTotal string
	StrMemoryFree string
	StrMemoryAvailable string
	StrSwapTotal string
	StrSwapFree string
	StrCPUSystem string
	StrCPUIdle string
	StrUpTime time.Duration
	StrDiskStats []disk.Stats
	StrNetworkStats []network.Stats
	StrLoadAvg1 string
	StrLoadAvg5 string
	StrLoadAvg15 string
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

var activeScreens ActiveScreens

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

		systemStats := collectSystemStats()
		log.Println(systemStats)
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

func collectSystemStats() PureSystemStat {
	memStats, err := memory.Get()
	if err != nil {
		log.Println(err)
	}

	cpuStats, err := cpu.Get()
	if err != nil {
		log.Println(err)
	}

	upTime, err := uptime.Get()
	if err != nil {
		log.Println(err)
	}

	diskStats, err := disk.Get()
	if err != nil {
		log.Println(err)
	}

	netStats, err := network.Get()
	if err != nil {
		log.Println(err)
	}

	loadAvg, err := loadavg.Get()
	if err != nil {
		log.Println(err)
	}

	log.Printf("MemoryTotal: %s\n", formatSizeUint64(memStats.Total))

	strSystemStat := StringSystemStat{
		StrMemoryTotal: formatSizeUint64(memStats.Total),
		StrMemoryFree: formatSizeUint64(memStats.Free),
		StrMemoryAvailable: formatSizeUint64(memStats.Available),
		StrSwapTotal: formatSizeUint64(memStats.SwapTotal),
		StrSwapFree: formatSizeUint64(memStats.SwapFree),
		StrCPUSystem: formatSizeUint64(cpuStats.System),
		StrCPUIdle: formatSizeUint64(cpuStats.Idle),
		StrLoadAvg1: formatSizeFloat64(loadAvg.Loadavg1),
		StrLoadAvg5: formatSizeFloat64(loadAvg.Loadavg5),
		StrLoadAvg15: formatSizeFloat64(loadAvg.Loadavg15),
	}

	log.Println(strSystemStat)

	return PureSystemStat{
		MemoryTotal: memStats.Total,
		MemoryFree: memStats.Free,
		MemoryAvailable: memStats.Available,
		SwapTotal: memStats.SwapTotal,
		SwapFree: memStats.SwapFree,
		CPUSystem: cpuStats.System,
		CPUIdle: cpuStats.Idle,
		UpTime: upTime,
		DiskStats: diskStats,
		NetworkStats: netStats,
		LoadAvg1: loadAvg.Loadavg1,
		LoadAvg5: loadAvg.Loadavg5,
		LoadAvg15: loadAvg.Loadavg15,
	}
}

func formatSizeUint64(data uint64) string {
	var units = [5]string{"B", "KB", "MB", "GB", "TB"}

	floatData := float64(data)

	i := 0
	for ; floatData > 1024; {
		floatData = floatData / 1024
		i++
	}

	s := strconv.FormatFloat(floatData, 'f', 6, 64) + units[i]
	return s
}

func formatSizeFloat64(data float64) string {
	s := strconv.FormatFloat(data, 'f', 3, 64)
	return s
}
