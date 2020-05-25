package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func BatteryInforamtion() {
	arg1 := "upower"
	arg2 := "-i"
	arg3 := "/org/freedesktop/UPower/devices/DisplayDevice"

	out, err := exec.Command(arg1, arg2, arg3).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}

	output := string(out[:])
	per := strings.Replace(output, "\n", ",", -1)
	per = strings.Replace(per, " ", "", -1)
	per = strings.TrimRight(per, ",")
	values := strings.Split(per, ",")

	hoursLeft := values[11][12:]
	BatteryPercentage := values[12][11:13]
	ChargingStatus := values[6][6:]
	Icon := values[13][10:]
	Icon = strings.Trim(Icon, "'")
	fmt.Println()
	fmt.Println("Remaining : ", hoursLeft)
	fmt.Println(BatteryPercentage)
	fmt.Println("Percentage : ", values[12][11:])
	fmt.Println(Icon)
	fmt.Println(ChargingStatus)
	fmt.Println()
	_, err = os.Stat("/home/root/Desktop/Text2Speech/status.txt")
	if err != nil {
		CheckBattery(BatteryPercentage, ChargingStatus, Icon)
	} else {
		return
	}

}
func CheckBattery(BatteryPercentage, ChargingStatus, Icon string) {
	Battery := strings.TrimRight(BatteryPercentage, ".")
	Percentage, _ := strconv.Atoi(Battery)
	if Percentage < 9 {
		message := fmt.Sprintf("Please connect the charger")
		Alert(message)
		Message := fmt.Sprintf("Hey Kuldeep, Battery is Low. Level: %s\n", BatteryPercentage)
		Title := "Battery Alert"
		Notification(BatteryPercentage, Title, Message, Icon)
	} else if Percentage == 99 {
		message := fmt.Sprintf("Charge full, Please Remove the charger")
		Alert(message)
		Message := fmt.Sprintf("Hey Kuldeep, Battery is Full. Level: %s\n", BatteryPercentage)
		Title := "Battery Alert"
		Notification(BatteryPercentage, Title, Message, Icon)
	}
}

func Notification(BatteryPercentage, Title, Message, Icon string) {
	notify := "notify-send"
	user := "critical"
	cmd := exec.Command(notify, Title, Message, "-u", user, "-i", Icon)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)

	}
}

func Alert(message string) {
	command := "espeak"
	voice := "-v"
	langauge := "en-us"
	gender := "+f3"
	voices := voice + langauge + gender
	speed := "-s130"
	cmd := exec.Command(command, voices, speed, message)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)

	}
	StatusFile()

}

func StatusFile() {
	// check if file exists
	path := "/home/root/Desktop/Text2Speech/status.txt"
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}
	fmt.Println("File Created Successfully", path)
}

func Speak() {
	// s := "Make the Computer speak"

	file, err := os.Open("data.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, speak := range txtlines {
		// fmt.Println(eachline)
		command := "espeak"
		voice := "-v"
		langauge := "en-us"
		gender := "+f3"
		voices := voice + langauge + gender
		speed := "-s130"
		cmd := exec.Command(command, voices, speed, speak)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)

		}
		time.Sleep(time.Second * 2)
	}
}

// Get the All info about wifi
func WifiStatus() {
	//  iwlist scan
	command := "sudo"
	arg := "iwlist"
	arg2 := "scan"
	out, err := exec.Command(command, arg, arg2).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	// fmt.Println(string(out))
	// address
	re := regexp.MustCompile(`Address: (.*)`)
	address := re.FindAllStringSubmatch(string(out), -1)
	// Name of Wifi
	re1 := regexp.MustCompile(`ESSID:"(.*|\n)"`)
	Name := re1.FindAllStringSubmatch(string(out), -1)
	// Bitrates
	re2 := regexp.MustCompile(`Bit Rates:(.*)`)
	Bits := re2.FindAllStringSubmatch(string(out), -1)

	// for i, addr := range address {
	for i := 0; i < len(address); i++ {
		fmt.Printf("Name %d : %s\n", i, Name[i][1])
		fmt.Printf("Address %d : %s\n", i, address[i][1])
		fmt.Printf("Bits %d : %s\n", i, Bits[i][1])
		fmt.Println("------------------------------------------------------------")
	}
	// }

	// space := strings.Replace(string(out), "\n", ",", -1)
	// space = strings.Replace(space, " ", "", -1)
	// // fmt.Println(space)
	// comma := strings.TrimRight(space, ",")
	// values := strings.Split(comma, ",")
	// fmt.Println(values)

}

func main() {
	var wg sync.WaitGroup
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
		os.Exit(130)
	} else {
		WifiStatus()
		os.Exit(130)
		// fmt.Println(runtime.GOOS)
		// Battery Percentage Alert Message
		ticker := time.NewTicker(5 * time.Second)
		quit := make(chan struct{})
		wg.Add(1)
		go func() {
			for {
				select {
				case <-ticker.C:
					BatteryInforamtion()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
		wg.Wait()
	}
}
