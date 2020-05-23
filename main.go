package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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
	fmt.Println()
	fmt.Println("Remaining : ", hoursLeft)
	fmt.Println("Percentage : ", values[12][11:])
	fmt.Println()
	_, err = os.Stat("status.txt")
	if err != nil {
		CheckBattery(BatteryPercentage, ChargingStatus)
	} else {
		return
	}

}
func CheckBattery(BatteryPercentage, ChargingStatus string) {
	if BatteryPercentage <= "10" && ChargingStatus == "discharging" {
		message := fmt.Sprintf("Please connect the charger")
		Alert(message)
	} else if BatteryPercentage == "99" && ChargingStatus == "charging" {
		message := fmt.Sprintf("Charge full, Please Remove the charger")
		Alert(message)
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
	path := "status.txt"
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

func main() {
	var wg sync.WaitGroup
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {
		// fmt.Println(runtime.GOOS)
		// Battery Percentage Alert Message
		ticker := time.NewTicker(5 * time.Minute)
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
