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

	var (
		tablesCount, costPerHour       int
		openingTimeStr, closingTimeStr string
	)

	scanner.Scan()
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &tablesCount); err != nil {
		fmt.Println(scanner.Text())
		return
	}

	scanner.Scan()
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
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &costPerHour); err != nil {
		fmt.Println(scanner.Text())
		return
	}

	computerClub := NewComputerClub(tablesCount, costPerHour, openingTime, closingTime)
	for scanner.Scan() { // TODO: check if time not sorted
		inputEvent := scanner.Text()
		if _, err := fmt.Fprintln(out, inputEvent); err != nil {
			log.Fatalln(err)
		}

		if outputEvent, err := computerClub.ProcessClientEvent(inputEvent); err != nil {
			out.Reset(nil)
			fmt.Println(inputEvent)
			return
		} else if _, err := fmt.Fprintln(out, outputEvent); err != nil {
			log.Fatalln(err)
		}
	}
}
