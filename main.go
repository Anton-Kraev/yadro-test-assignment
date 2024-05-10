package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
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
		placesCount, costPerHour       int
		openingTimeStr, closingTimeStr string
	)

	scanner.Scan()
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &placesCount); err != nil {
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

	computerClub := NewComputerClub(placesCount, costPerHour, openingTime, closingTime)
	if _, err := fmt.Fprintf(out, "%02d:%02d\n", openingTime.Hour(), openingTime.Minute()); err != nil {
		log.Fatalln(err)
	}
	for scanner.Scan() {
		inputEvent := scanner.Text()
		if _, err := fmt.Fprintln(out, inputEvent); err != nil {
			log.Fatalln(err)
		}

		outputEvent, err := computerClub.ProcessClientEvent(inputEvent)
		if err != nil {
			out.Reset(nil)
			fmt.Println(inputEvent)
			return
		}
		if outputEvent != "" {
			if _, err := fmt.Fprintln(out, outputEvent); err != nil {
				log.Fatalln(err)
			}
		}
	}
	for _, client := range computerClub.Close() {
		if _, err := fmt.Fprintf(out, "%02d:%02d 11 %s\n", closingTime.Hour(), closingTime.Minute(), client); err != nil {
			log.Fatalln(err)
		}
	}
	if _, err := fmt.Fprintf(out, "%02d:%02d\n", closingTime.Hour(), closingTime.Minute()); err != nil {
		log.Fatalln(err)
	}
	sort.Slice(computerClub.placeStats, func(i, j int) bool {
		return computerClub.placeStats[i].id < computerClub.placeStats[j].id
	})
	for _, place := range computerClub.placeStats {
		if _, err := fmt.Fprintf(out, "%d %d %02d:%02d\n", place.id, place.revenue*costPerHour, place.useTime.Hour(), place.useTime.Minute()); err != nil {
			log.Fatalln(err)
		}
	}
}
