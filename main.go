package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	file, err := os.Open("input.txt")
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}(file)
	if err != nil {
		log.Fatalln(err)
	}
	scanner := bufio.NewScanner(file)
	out := bufio.NewWriter(os.Stdout)
	defer func(out *bufio.Writer) {
		if err := out.Flush(); err != nil {
			log.Fatalln(err)
		}
	}(out)

	scanner.Scan()
	var tablesCount int
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &tablesCount); err != nil {
		fmt.Println(scanner.Text())
		return
	}

	scanner.Scan()
	var openingTimeStr, closingTimeStr string
	if _, err := fmt.Sscanf(scanner.Text(), "%s %s", &openingTimeStr, &closingTimeStr); err != nil {
		fmt.Println(scanner.Text())
		return
	}
	openingTime, err1 := time.Parse("15:04", openingTimeStr)
	closingTime, err2 := time.Parse("15:04", closingTimeStr)
	if errors.Join(err1, err2) != nil || !closingTime.After(openingTime) {
		fmt.Println(scanner.Text())
		return
	}

	scanner.Scan()
	var costPerHour int
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &costPerHour); err != nil {
		fmt.Println(scanner.Text())
		return
	}

	for scanner.Scan() {
		inputEvent := scanner.Text()
		_, err := fmt.Fprintln(out, inputEvent)
		if err != nil {
			log.Fatalln(err)
		}
		//outputEvents, err := 0, new(error)
	}
	// out.Reset(os.Stdout)
}
