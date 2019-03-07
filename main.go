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

/* string os stats structs */
type StrCPUStat struct {
	User string `json:"user"`
	System string `json:"system"`
	Idle string `json:"idle"`
}

type StrLoadAverage struct {
	Avg1 string `json:"avg1"`
	Avg5 string `json:"avg5"`
	Avg15 string `json:"avg15"`
}

type StrOsStat struct {
	Memory StrMemoryStat `json:"memoryStat"`
	CPU StrCPUStat `json:"cpuStat"`
	Uptime string `json:"uptime"`
	Disk []StrDiskStats `json:"diskStat"`
	Network []StrNetStats `json:"networkStat"`
	LoadAverage StrLoadAverage `json:"loadAverage"`
}

type StrMemoryStat struct {
	Total string `json:"total"`
	Free string `json:"free"`
	Available string `json:"available"`
	SwapTotal string `json:"swapTotal"`
	SwapFree string `json:"swapFree"`
}

type StrDiskStats struct {
	Name string `json:"name"`
	ReadsCompleted string `json:"readsCompleted"`
	WritesCompleted string `json:"writesCompleted"`
}

type StrNetStats struct {
	Name string `json:"name"`
	RxBytes string `json:"rxBytes"`
	TxBytes string `json:"txBytes"`
}

/* pure os stat structs */

type OsStat struct {
	Memory MemoryStat `json:"memoryStat"`
	CPU CPUStat `json:"cpuStat"`
	Uptime time.Duration `json:"uptime"`
	Disk []disk.Stats `json:"diskStat"`
	Network []network.Stats `json:"networkStat"`
	LoadAvg LoadAverage `json:"loadAverage"`
}

type MemoryStat struct {
	Total uint64 `json:"total"`
	Free uint64 `json:"free"`
	Available uint64 `json:"available"`
	SwapTotal uint64 `json:"swapTotal"`
	SwapFree uint64 `json:"swapFree"`
}

type CPUStat struct {
	User uint64 `json:"user"`
	System uint64 `json:"system"`
	Idle uint64 `json:"idle"`
}

type LoadAverage struct {
	Avg1 float64 `json:"avg1"`
	Avg5 float64 `json:"avg5"`
	Avg15 float64 `json:"avg15"`
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

		systemStats, strSystemStats := collectSystemStats()
		log.Println(systemStats)
		log.Println(strSystemStats)
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

func collectSystemStats() (OsStat, StrOsStat) {
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

	/* parse disk stats */
	var strDiskStats []StrDiskStats

	for i := 0; i <= len(diskStats)-1; i++ {
		sReadsCompleted := formatSizeUint64(diskStats[i].ReadsCompleted)
		sWritesCompleted := formatSizeUint64(diskStats[i].WritesCompleted)

		tempDiskStats := StrDiskStats {
			Name: diskStats[i].Name,
			ReadsCompleted: sReadsCompleted,
			WritesCompleted: sWritesCompleted,
		}
		strDiskStats = append(strDiskStats, tempDiskStats)
	}

	/* parse network stats */
	var strNetStats []StrNetStats

	for i := 0; i <= len(netStats)-1; i++ {
		sRxBytes := formatSizeUint64(netStats[i].RxBytes)
		sTxBytes := formatSizeUint64(netStats[i].TxBytes)

		tempNetStats := StrNetStats {
			Name: netStats[i].Name,
			RxBytes: sRxBytes,
			TxBytes: sTxBytes,
		}
		strNetStats = append(strNetStats, tempNetStats)
	}

	//upTimeStr := upTime.String()

	//log.Println(time.Parse(time.UnixDate, upTimeStr))

	//log.Println(upTime.String())
	//log.Println(upTime.Format("2006-01-02 15:04:05"))
	//log.Println(strSystemStat)


	return OsStat {
		Memory: MemoryStat {
			Total: memStats.Total,
			Free: memStats.Free,
			Available: memStats.Available,
			SwapTotal: memStats.SwapTotal,
			SwapFree: memStats.SwapFree,
		},
		CPU: CPUStat {
			User: cpuStats.User,
			System: cpuStats.System,
			Idle: cpuStats.Idle,
		},
		Uptime: upTime,
		Disk: diskStats,
		Network: netStats,
		LoadAvg: LoadAverage{
			Avg1: loadAvg.Loadavg1,
			Avg5: loadAvg.Loadavg5,
			Avg15: loadAvg.Loadavg15,
		},
	}, StrOsStat {
		Memory: StrMemoryStat {
			Total: formatSizeUint64(memStats.Total),
			Free: formatSizeUint64(memStats.Free),
			Available: formatSizeUint64(memStats.Available),
			SwapTotal: formatSizeUint64(memStats.SwapTotal),
			SwapFree: formatSizeUint64(memStats.SwapFree),
		},
		CPU: StrCPUStat {
			User: strconv.FormatUint(cpuStats.System, 10),
			System: strconv.FormatUint(cpuStats.System, 10),
			Idle: strconv.FormatUint(cpuStats.Idle, 10),
		},
		Uptime: upTime.String(),
		Disk: strDiskStats,
		Network: strNetStats,
		LoadAverage: StrLoadAverage {
			Avg1: formatSizeFloat64(loadAvg.Loadavg1),
			Avg5: formatSizeFloat64(loadAvg.Loadavg5),
			Avg15: formatSizeFloat64(loadAvg.Loadavg15),
		},
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

	s := strconv.FormatFloat(floatData, 'f', 2, 64) + units[i]
	return s
}

func formatSizeFloat64(data float64) string {
	s := strconv.FormatFloat(data, 'f', 2, 64)
	return s
}
