package stat

import(
	"strconv"
	"time"
	"log"
	"os"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"

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
	Percentage string `json:"percentage"`
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
	Percentage float64 `json:"percentage"`
}

type LoadAverage struct {
	Avg1 float64 `json:"avg1"`
	Avg5 float64 `json:"avg5"`
	Avg15 float64 `json:"avg15"`
}

/* Holds screen informations from command */
type SystemScreen struct {
	Name string `json:"name"`
	Up bool `json:"up"`
}

type SystemScreens struct {
	Screens []SystemScreen `json:"screens"`
}

/* Holds name of the screen that
 * we want to check whether active or not */
type ActiveScreens struct {
	Names []string `json:"activeScreen"`
}

var activeScreens ActiveScreens

func CollectSystemStats() (*OsStat, error) {
	memStats, err := memory.Get()
	if err != nil {
		return nil, err
	}

	cpuStats, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	upTime, err := uptime.Get()
	if err != nil {
		return nil, err
	}

	diskStats, err := disk.Get()
	if err != nil {
		return nil, err
	}

	netStats, err := network.Get()
	if err != nil {
		return nil, err
	}

	loadAvg, err := loadavg.Get()
	if err != nil {
		return nil, err
	}

	//upTimeStr := upTime.String()

	//log.Println(time.Parse(time.UnixDate, upTimeStr))

	//log.Println(upTime.String())
	//log.Println(upTime.Format("2006-01-02 15:04:05"))
	//log.Println(strSystemStat)

	cpuPercentage := float64((cpuStats.User + cpuStats.System)) / float64(cpuStats.Idle)
	cpuPercentage = cpuPercentage * 100
	log.Println("cpu percentage", cpuPercentage)

	return &OsStat {
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
			Percentage: cpuPercentage,
		},
		Uptime: upTime,
		Disk: diskStats,
		Network: netStats,
		LoadAvg: LoadAverage{
			Avg1: loadAvg.Loadavg1,
			Avg5: loadAvg.Loadavg5,
			Avg15: loadAvg.Loadavg15,
		},
	}, nil
}

func CollectStrSystemStats() (*StrOsStat, error) {
	osStats, err := CollectSystemStats()
	if err != nil {
		return nil, err
	}

	log.Println("osstats")
	log.Println(osStats)

		/* parse disk stats */
	var strDiskStats []StrDiskStats

	for i := 0; i <= len(osStats.Disk)-1; i++ {
		sReadsCompleted := formatSizeUint64(osStats.Disk[i].ReadsCompleted)
		sWritesCompleted := formatSizeUint64(osStats.Disk[i].WritesCompleted)

		tempDiskStats := StrDiskStats {
			Name: osStats.Disk[i].Name,
			ReadsCompleted: sReadsCompleted,
			WritesCompleted: sWritesCompleted,
		}
		strDiskStats = append(strDiskStats, tempDiskStats)
	}

	/* parse network stats */
	var strNetStats []StrNetStats

	for i := 0; i <= len(osStats.Network)-1; i++ {
		sRxBytes := formatSizeUint64(osStats.Network[i].RxBytes)
		sTxBytes := formatSizeUint64(osStats.Network[i].TxBytes)

		log.Println("osstats network")
		log.Println(osStats.Network[i])
		tempNetStats := StrNetStats {
			Name: osStats.Disk[i].Name,
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

	return &StrOsStat {
		Memory: StrMemoryStat {
			Total: formatSizeUint64(osStats.Memory.Total),
			Free: formatSizeUint64(osStats.Memory.Free),
			Available: formatSizeUint64(osStats.Memory.Available),
			SwapTotal: formatSizeUint64(osStats.Memory.SwapTotal),
			SwapFree: formatSizeUint64(osStats.Memory.SwapFree),
		},
		CPU: StrCPUStat {
			User: strconv.FormatUint(osStats.CPU.User, 10),
			System: strconv.FormatUint(osStats.CPU.System, 10),
			Idle: strconv.FormatUint(osStats.CPU.Idle, 10),
			Percentage: "%"+formatSizeFloat64(osStats.CPU.Percentage),
		},
		Uptime: osStats.Uptime.String(),
		Disk: strDiskStats,
		Network: strNetStats,
		LoadAverage: StrLoadAverage {
			Avg1: formatSizeFloat64(osStats.LoadAvg.Avg1),
			Avg5: formatSizeFloat64(osStats.LoadAvg.Avg5),
			Avg15: formatSizeFloat64(osStats.LoadAvg.Avg15),
		},
	}, nil
}

func CollectScreenStats() (*SystemScreens, error) {
	activeScreen, err := readActiveScreensConfig("active_screen.json")
	if err != nil {
		return nil, err
	}

	systemScreen := updateSystemScreen()
	checkScreens := CheckScreens(activeScreen, systemScreen)

	return &SystemScreens{
		Screens: checkScreens,
	}, nil
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

func CheckScreens(activeScreens *ActiveScreens, systemScreens []string) []SystemScreen{
	var sysScreens []SystemScreen

	for _, value := range activeScreens.Names {
		var systemScreen SystemScreen
		exists := contains(systemScreens, value)
		if exists {
			systemScreen = SystemScreen{
				Name: value,
				Up: true,
			}
			log.Printf("%s is running.\n", value)
		} else {
			systemScreen = SystemScreen{
				Name: value,
				Up: false,
			}
			sendCrashMessage(value)
			log.Printf("%s is not running\n", value)
		}
		sysScreens = append(sysScreens, systemScreen)
	}

	return sysScreens
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
	var sysScreenNames []string

	out, err := exec.Command("screen", "-ls").Output()
	if err != nil {
		log.Printf("Error while executing screen -ls command: %s\n", err)
		return sysScreenNames
	}
	screenOut := string(out)
	screenOutLineArr := strings.Split(screenOut, "\n")

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

func readActiveScreensConfig(configFileName string) (*ActiveScreens, error) {
	var activeScreen ActiveScreens
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Reading Active Screens: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&activeScreen); err != nil {
		return nil, err
	}

	return &activeScreen, nil
}
