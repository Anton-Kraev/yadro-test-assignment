package main

import (
	"errors"
	"strconv"
	"strings"
	t "time"
)

type ComputerClub struct {
	tablesCount, costPerHour                int
	openingTime, closingTime, prevEventTime t.Time
}

func NewComputerClub(tablesCount, costPerHour int, openingTime, closingTime t.Time) *ComputerClub {
	return nil
}

func (club *ComputerClub) ProcessClientEvent(event string) (string, error) {
	processingError := errors.New("bad event format")

	eventInfo := strings.Split(event, " ")
	if len(eventInfo) != 3 && len(eventInfo) != 4 {
		return "", processingError
	}

	var (
		time                   t.Time
		eventType, tableNumber int
		clientName             string
		err                    error
	)

	if time, err = t.Parse("15:04", eventInfo[0]); err != nil || club.prevEventTime.After(time) {
		return "", processingError
	}
	club.prevEventTime = time

	if eventType, err = strconv.Atoi(eventInfo[1]); err != nil {
		return "", processingError
	}

	if clientName = eventInfo[2]; len(clientName) == 0 {
		return "", processingError
	}

	switch eventType {
	case 1:
		return club.clientCame(time, clientName)
	case 2:
		if len(eventInfo) != 4 {
			return "", processingError
		}
		if tableNumber, err = strconv.Atoi(eventInfo[3]); tableNumber <= 0 || tableNumber > club.tablesCount {
			return "", processingError
		}
		return club.clientSat(time, clientName, tableNumber)
	case 3:
		return club.clientWaiting(time, clientName)
	case 4:
		return club.clientLeft(time, clientName)
	}

	return "", processingError
}

func (club *ComputerClub) clientCame(time t.Time, clientName string) (string, error) {
	return "", errors.New("")
}

func (club *ComputerClub) clientSat(time t.Time, clientName string, tableNumber int) (string, error) {
	return "", errors.New("")
}

func (club *ComputerClub) clientWaiting(time t.Time, clientName string) (string, error) {
	return "", errors.New("")
}

func (club *ComputerClub) clientLeft(time t.Time, clientName string) (string, error) {
	return "", errors.New("")
}
