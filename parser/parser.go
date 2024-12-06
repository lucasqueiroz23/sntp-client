package parser

import (
	"bytes"
	"encoding/binary"
	// "fmt"
	"sntp-client/client-socket"
	"sntp-client/error-handling"
	"strconv"
)

const baseYear uint32 = 1900       // int NTP, the base date is jan 1, 1900
const daysInAYear float64 = 365.25 // the .25 accounts for years that have 366 days
const hoursInADay float64 = 24.0
const minutesInAnHour float64 = 60.0
const secondsInAMinute float64 = 60.0

const dateTimePrefix string = "Data/hora: "

func GetDate(serverResponse []byte) string {
	timePassedInSeconds := getResponsePacket(serverResponse).TxTm_s
	currentYear := getCurrentYear(timePassedInSeconds)
	today := getDayOfTheWeek(daysPastSince(timePassedInSeconds))

	return dateTimePrefix +
		strconv.FormatUint(uint64(currentYear), 10) + " " +
		today
}

func getResponsePacket(serverResponse []byte) *clientSocket.NtpPacket {
	response := new(clientSocket.NtpPacket)

	binaryResponse := bytes.NewReader(serverResponse)

	if readErr := binary.Read(binaryResponse, binary.BigEndian, response); readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	return response
}

func getCurrentYear(timePassedInSeconds uint32) uint32 {
	return baseYear + uint32(float64(timePassedInSeconds)/(daysInAYear*hoursInADay*minutesInAnHour*secondsInAMinute))
}

func daysPastSince(timePassedInSeconds uint32) uint32 {
	return timePassedInSeconds / (uint32(hoursInADay) * uint32(minutesInAnHour) * uint32(secondsInAMinute))
}

func getDayOfTheWeek(daysSinceBaseDate uint32) string {
	var daysOfTheWeek = [7]string{"Seg", "Ter", "Qua", "Qui", "Sex", "Sáb", "Dom"}
	const daysInAWeek = 7

	return daysOfTheWeek[daysSinceBaseDate%daysInAWeek]
}
