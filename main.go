package main

import(
	"os/exec"
	"fmt"
	"strings"
	"log"
	//"regexp"
)

func main() {
	out, err := exec.Command("screen", "-ls").Output()

	//xregexString, _ := regexp.Compile("p([a-z]+)ch")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(out)
	screenOut := string(out)
	//screenOut = strings.TrimSpace(screenOut)
	screenOutArr := strings.Split(screenOut, "\n")
	screenOut = screenOutArr[1]
	screenOut = strings.TrimSpace(screenOut)
	screenOutArr = strings.Split(screenOut, "\t")
	screenOut = screenOutArr[0]
	fmt.Println(screenOut)

	screenOutArr = strings.Split(screenOut, ".")
	screenPID := screenOutArr[0]
	screenName := screenOutArr[1]

	fmt.Println("screenPID:" + screenPID)
	fmt.Println("screenName:" + screenName)
}
