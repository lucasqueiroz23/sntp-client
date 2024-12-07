package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sntp-client/client-socket"
	"sntp-client/error-handling"
	"strconv"
)

const baseYear uint32 = 1900       // int NTP, the base date is jan 1, 1900
const daysInAYear float64 = 365.25 // the .25 accounts for years that have 366 days
const daysInAMonth float64 = 30.44 // the .44 account for leap years aswell
const hoursInADay float64 = 24.0
const minutesInAnHour float64 = 60.0
const secondsInAMinute float64 = 60.0

const dateTimePrefix string = "Data/hora: "

func GetDate(serverResponse []byte) string {
	timePassedInSeconds := getResponsePacket(serverResponse).TxTm_s
	currentYear := getCurrentYear(timePassedInSeconds)

	daysSinceBaseDate := daysPastSince(timePassedInSeconds)
	monthsSinceBaseDate := monthsPastSince(timePassedInSeconds)

	today := getDayOfTheWeek(daysSinceBaseDate)
	currentMonth := getCurrentMonth(monthsSinceBaseDate)

	currentTime := getCurrentTime(timePassedInSeconds)

	currentDay := "not implemented yet"

	return dateTimePrefix +
		today + " " +
		currentMonth + " " +
		currentDay + " " +
		currentTime + " " +
		strconv.FormatUint(uint64(currentYear), 10)
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

func monthsPastSince(timePassedInSeconds uint32) uint32 {
	return uint32(float64(timePassedInSeconds) / ((daysInAMonth) * (hoursInADay) * (minutesInAnHour) * (secondsInAMinute)))
}

func getCurrentMonth(monthsPassedSinceBaseDate uint32) string {
	var months = [12]string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
	const monthsInAYear = 12
	return months[monthsPassedSinceBaseDate%monthsInAYear]
}

func getCurrentTime(timePassedInSeconds uint32) string {
	seconds := timePassedInSeconds % 60
	minutes := (timePassedInSeconds / 60) % 60
	hours := (timePassedInSeconds/(60*60) - 3) % 24
	fmt.Println(timePassedInSeconds)
	return strconv.FormatInt(int64(hours), 10) + ":" + strconv.FormatInt(int64(minutes), 10) + ":" + strconv.FormatInt(int64(seconds), 10)
}

// func getCurrentDay(currentMonth string) string {
//
// 	daysInMonths := map[string]int{
// 		"Jan": 31,
// 		"Fev": 28,
// 		"Mar": 31,
// 		"Abr": 30,
// 		"Mai": 31,
// 		"Jun": 30,
// 		"Jul": 31,
// 		"Ago": 31,
// 		"Set": 30,
// 		"Out": 31,
// 		"Nov": 30,
// 		"Dez": 31,
// 	}
//
// 	return strconv.FormatInt(int64(daysInMonths[currentMonth]), 10)
// }
