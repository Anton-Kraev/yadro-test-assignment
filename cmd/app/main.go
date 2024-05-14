package main

import (
	"log"
	"os"
	computerclub "yadro-test-assignment/internal/computer_club"
)

func main() {
	file, err := os.Open(os.Args[1])
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}(file)
	if err != nil {
		log.Fatalln(err)
	}
	computerclub.ProcessComputerClubDayEvents(file, os.Stdout)
}
