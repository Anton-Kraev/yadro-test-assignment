package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	t "time"
)

type ComputerPlace struct {
	id, revenue           int
	occupiedTime, useTime t.Time
	occupied              bool
}

type ComputerClub struct {
	placesCount, costPerHour                int
	openingTime, closingTime, prevEventTime t.Time
	placeStats                              []ComputerPlace
	clientPlace                             map[string]int
	queue                                   []string
}

func NewComputerClub(placesCount, costPerHour int, openingTime, closingTime t.Time) *ComputerClub {
	prevEventTime, _ := t.Parse("15:04", "00:00")
	placeStats := make([]ComputerPlace, placesCount)
	zeroTime, _ := t.Parse("15:04", "00:00")
	for i := 0; i < placesCount; i++ {
		placeStats[i] = ComputerPlace{id: i + 1, occupiedTime: zeroTime, useTime: zeroTime}
	}
	return &ComputerClub{
		placesCount:   placesCount,
		costPerHour:   costPerHour,
		openingTime:   openingTime,
		closingTime:   closingTime,
		prevEventTime: prevEventTime,
		placeStats:    placeStats,
		clientPlace:   make(map[string]int),
		queue:         make([]string, 0),
	}
}

func (club *ComputerClub) AcquirePlace(place int, client string, time t.Time) {
	club.clientPlace[client] = place
	club.placeStats[place-1].occupied = true
	club.placeStats[place-1].occupiedTime = time
}

func (club *ComputerClub) ReleasePlace(place int, time t.Time) {
	if place == 0 {
		return
	}
	placeStats := &club.placeStats[place-1]
	placeStats.useTime = placeStats.useTime.Add(time.Sub(placeStats.occupiedTime))
	placeStats.revenue += int(math.Ceil(time.Sub(placeStats.occupiedTime).Hours()))
	placeStats.occupiedTime, _ = t.Parse("15:04", "00:00")
	placeStats.occupied = false
}

func (club *ComputerClub) ReleaseClientPlace(client string, time t.Time) {
	place := club.clientPlace[client]
	delete(club.clientPlace, client)
	club.ReleasePlace(place, time)
}

func (club *ComputerClub) ProcessClientEvent(event string) (string, error) {
	processingError := errors.New("bad event format")

	eventInfo := strings.Split(event, " ")
	if len(eventInfo) != 3 && len(eventInfo) != 4 {
		return "", processingError
	}

	var (
		time                   t.Time
		eventType, placeNumber int
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
		if placeNumber, err = strconv.Atoi(eventInfo[3]); placeNumber <= 0 || placeNumber > club.placesCount {
			return "", processingError
		}
		return club.clientSat(time, clientName, placeNumber)
	case 3:
		return club.clientWaiting(time, clientName)
	case 4:
		return club.clientLeft(time, clientName)
	}

	return "", processingError
}

func (club *ComputerClub) Close() []string {
	var clientsInClub []string
	for client := range club.clientPlace {
		clientsInClub = append(clientsInClub, client)
		club.ReleaseClientPlace(client, club.closingTime)
	}
	sort.Strings(clientsInClub)
	return clientsInClub
}

func (club *ComputerClub) clientCame(time t.Time, clientName string) (string, error) {
	if time.After(club.closingTime) {
		return "", errors.New("failed to process the clientCame event")
	}

	_, inClub := club.clientPlace[clientName]
	if inClub {
		return fmt.Sprintf("%02d:%02d 13 YouShallNotPass", time.Hour(), time.Minute()), nil
	} else if time.Before(club.openingTime) {
		return fmt.Sprintf("%02d:%02d 13 NotOpenYet", time.Hour(), time.Minute()), nil
	}
	club.clientPlace[clientName] = 0
	return "", nil
}

func (club *ComputerClub) clientSat(time t.Time, clientName string, placeNumber int) (string, error) {
	placeId, inClub := club.clientPlace[clientName]
	if !inClub || time.Before(club.openingTime) || time.After(club.closingTime) {
		return "", errors.New("failed to process the clientSat event")
	}

	if club.placeStats[placeNumber-1].occupied {
		return fmt.Sprintf("%02d:%02d 13 PlaceIsBusy", time.Hour(), time.Minute()), nil
	}
	club.ReleasePlace(placeId, time)
	club.AcquirePlace(placeNumber, clientName, time)
	return "", nil
}

func (club *ComputerClub) clientWaiting(time t.Time, clientName string) (string, error) {
	placeId, inClub := club.clientPlace[clientName]
	if !inClub || placeId != 0 || time.Before(club.openingTime) || time.After(club.closingTime) {
		return "", errors.New("failed to process the clientWaiting event")
	}

	for _, place := range club.placeStats {
		if !place.occupied {
			return fmt.Sprintf("%02d:%02d 13 ICanWaitNoLonger!", time.Hour(), time.Minute()), nil
		}
	}
	if len(club.queue) >= len(club.placeStats) {
		return fmt.Sprintf("%02d:%02d 11 %s", time.Hour(), time.Minute(), clientName), nil
	}
	club.queue = append(club.queue, clientName)
	return "", nil
}

func (club *ComputerClub) clientLeft(time t.Time, clientName string) (string, error) {
	if time.Before(club.openingTime) || time.After(club.closingTime) {
		return "", errors.New("failed to process the clientLeft event")
	}

	placeId, inClub := club.clientPlace[clientName]
	if !inClub {
		return fmt.Sprintf("%02d:%02d 13 ClientUnknown", time.Hour(), time.Minute()), nil
	}
	club.ReleaseClientPlace(clientName, time)
	if len(club.queue) > 0 {
		clientInQueue := club.queue[0]
		club.queue = club.queue[1:]
		club.AcquirePlace(placeId, clientInQueue, time)
		return fmt.Sprintf("%02d:%02d 12 %s %d", time.Hour(), time.Minute(), clientInQueue, placeId), nil
	}
	return "", nil
}
