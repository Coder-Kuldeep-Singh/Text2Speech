package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
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
	// s := "Make the Computer speak"

}
