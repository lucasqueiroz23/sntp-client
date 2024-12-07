package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sntp-client/client-socket"
	"sntp-client/error-handling"
	"strconv"
)

const (
	baseTen                  = 10
	baseYear         uint32  = 1900   // in NTP, the base date is jan 1, 1900
	daysInAYear      float64 = 365.25 // the .25 accounts for years that have 366 days
	daysInAMonth     float64 = 30.44  // the .44 account for leap years aswell
	hoursInADay      float64 = 24.0
	minutesInAnHour  float64 = 60.0
	secondsInAMinute float64 = 60.0
	daysInAWeek              = 7
	monthsInAYear            = 12
)

var (
	daysOfTheWeek = [7]string{"Seg", "Ter", "Qua", "Qui", "Sex", "Sáb", "Dom"}
	months        = [12]string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
)

type monthData struct {
	daysPassedInTheYear, daysInTheMonth int
}

const dateTimePrefix string = "Data/hora: "

func GetDate(serverResponse []byte) string {
	timePassedInSeconds := getResponsePacket(serverResponse).TxTm_s
	currentYear := getCurrentYear(timePassedInSeconds)

	daysSinceBaseDate := daysPastSince(timePassedInSeconds)

	today := getDayOfTheWeek(daysSinceBaseDate)

	currentTime := getCurrentTime(timePassedInSeconds)

	currentDay := getCurrentDay(daysSinceBaseDate, currentYear)

	return dateTimePrefix +
		today + " " +
		currentDay + " " +
		currentTime + " " +
		fmt.Sprintf("%d", currentYear)
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
	return daysOfTheWeek[daysSinceBaseDate%daysInAWeek]
}

func monthsPastSince(timePassedInSeconds uint32) uint32 {
	return uint32(float64(timePassedInSeconds) / ((daysInAMonth) * (hoursInADay) * (minutesInAnHour) * (secondsInAMinute)))
}

func getCurrentTime(timePassedInSeconds uint32) string {

	seconds := timePassedInSeconds % uint32(secondsInAMinute)
	minutes := (timePassedInSeconds / uint32(secondsInAMinute)) % uint32(minutesInAnHour)
	hours := (timePassedInSeconds / (uint32(secondsInAMinute) * uint32(minutesInAnHour))) % uint32(hoursInADay)

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getCurrentDay(daysSinceBaseDate uint32, currentYear uint32) string {

	yearsPassed := currentYear - baseYear
	numberOfLeapYears := countLeapYears(int(baseYear), int(currentYear))

	daysPassedThisYear := daysSinceBaseDate - (365*(yearsPassed-uint32(numberOfLeapYears)) + 366*uint32(numberOfLeapYears))
	daysInMonths := [12]uint32{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	if yearIsLeap(currentYear) {
		daysInMonths[1]++
	}

	// check if feb 29 has passed
	if daysPassedThisYear > uint32(daysInMonths[0]+daysInMonths[1]) {
		daysPassedThisYear++
	}

	var currentDay uint32 = 0

	var currentMonthIndex int = 0
	for i := 0; i < 11; i++ {
		daysInMonths[i+1] += daysInMonths[i]

		if daysPassedThisYear >= daysInMonths[i] && daysPassedThisYear <= daysInMonths[i+1] {
			currentMonthIndex = i
			offsetToday := 1
			currentDay = (daysPassedThisYear - daysInMonths[currentMonthIndex]) + uint32(offsetToday)

			if currentDay > 0 {
				currentMonthIndex++
			}
		}
	}

	return months[currentMonthIndex] + " " + strconv.FormatUint(uint64(currentDay), baseTen)
}

func leapYearsUpTo(year int) int {
	return (year / 4) - (year / 100) + (year / 400)
}

func countLeapYears(start, end int) int {
	return leapYearsUpTo(end) - leapYearsUpTo(start-1)
}

func yearIsLeap(year uint32) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}
