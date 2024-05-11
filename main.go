package main

import (
	"log"
	"os"
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
	ProcessComputerClubDayEvents(file, os.Stdout)
}
