package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Battery Percentage Alert Message
	arg1 := "upower"
	arg2 := "-i"
	arg3 := "/org/freedesktop/UPower/devices/DisplayDevice"

	out, err := exec.Command(arg1, arg2, arg3).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Println(output)
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
		cmd := exec.Command("espeak", speak)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)

		}
		time.Sleep(time.Second * 2)
	}

}
