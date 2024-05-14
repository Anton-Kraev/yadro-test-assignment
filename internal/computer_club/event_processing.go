package computer_club

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	t "time"
)

func checkFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func tryPrintFatal(_ any, err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ProcessComputerClubDayEvents(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	out := bufio.NewWriter(w)
	defer func(out *bufio.Writer) {
		checkFatal(out.Flush())
	}(out)

	var (
		placesCount, costPerHour       int
		openingTimeStr, closingTimeStr string
	)

	scanner.Scan()
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &placesCount); err != nil || placesCount <= 0 {
		tryPrintFatal(fmt.Fprintln(w, scanner.Text()))
		return
	}

	scanner.Scan()
	if _, err := fmt.Sscanf(scanner.Text(), "%s %s", &openingTimeStr, &closingTimeStr); err != nil {
		tryPrintFatal(fmt.Fprintln(w, scanner.Text()))
		return
	}
	openingTime, err1 := t.Parse("15:04", openingTimeStr)
	closingTime, err2 := t.Parse("15:04", closingTimeStr)
	if errors.Join(err1, err2) != nil || !closingTime.After(openingTime) {
		tryPrintFatal(fmt.Fprintln(w, scanner.Text()))
		return
	}

	scanner.Scan()
	if _, err := fmt.Sscanf(scanner.Text(), "%d", &costPerHour); err != nil || costPerHour <= 0 {
		tryPrintFatal(fmt.Fprintln(w, scanner.Text()))
		return
	}

	computerClub := NewComputerClub(placesCount, costPerHour, openingTime, closingTime)

	tryPrintFatal(fmt.Fprintf(out, "%02d:%02d\n", openingTime.Hour(), openingTime.Minute()))
	for scanner.Scan() {
		inputEvent := scanner.Text()
		tryPrintFatal(fmt.Fprintln(out, inputEvent))
		outputEvent, err := processEvent(computerClub, inputEvent)
		if err != nil {
			out.Reset(nil)
			tryPrintFatal(fmt.Fprintln(w, inputEvent))
			return
		} else if outputEvent != "" {
			tryPrintFatal(fmt.Fprintln(out, outputEvent))
		}
	}
	for _, client := range computerClub.Close() {
		tryPrintFatal(fmt.Fprintf(out, "%02d:%02d 11 %s\n", closingTime.Hour(), closingTime.Minute(), client))
	}
	tryPrintFatal(fmt.Fprintf(out, "%02d:%02d\n", closingTime.Hour(), closingTime.Minute()))

	sort.Slice(computerClub.placeStats, func(i, j int) bool {
		return computerClub.placeStats[i].id < computerClub.placeStats[j].id
	})
	for _, place := range computerClub.placeStats {
		tryPrintFatal(fmt.Fprintf(
			out, "%d %d %02d:%02d\n",
			place.id, place.revenue*costPerHour, place.useTime.Hour(), place.useTime.Minute(),
		))
	}
}

func processEvent(club *ComputerClub, event string) (string, error) {
	processingError := errors.New("bad event format")

	var (
		time                   t.Time
		eventType, placeNumber int
		clientName             string
		err                    error
	)

	eventInfo := strings.Split(event, " ")
	if len(eventInfo) != 3 && len(eventInfo) != 4 {
		return "", processingError
	}
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
		return club.ClientCame(time, clientName)
	case 2:
		if len(eventInfo) != 4 {
			return "", processingError
		}
		if placeNumber, err = strconv.Atoi(eventInfo[3]); placeNumber <= 0 || placeNumber > club.placesCount {
			return "", processingError
		}
		return club.ClientSat(time, clientName, placeNumber)
	case 3:
		return club.ClientWaiting(time, clientName)
	case 4:
		return club.ClientLeft(time, clientName)
	}

	return "", processingError
}
