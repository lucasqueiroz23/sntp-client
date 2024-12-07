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

func getCurrentTime(timePassedInSeconds uint32) string {

	seconds := timePassedInSeconds % uint32(secondsInAMinute)
	minutes := (timePassedInSeconds / uint32(secondsInAMinute)) % uint32(minutesInAnHour)
	hours := (timePassedInSeconds / (uint32(secondsInAMinute) * uint32(minutesInAnHour))) % uint32(hoursInADay)

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getCurrentDay(daysSinceBaseDate uint32, currentYear uint32) string {
	daysInMonths := [12]uint32{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	yearsPassed := currentYear - baseYear
	numberOfLeapYears := countLeapYears(int(baseYear), int(currentYear))

	// days passed since jan 1st
	daysPassedSinceJanFirst := daysSinceBaseDate - (365*(yearsPassed-uint32(numberOfLeapYears)) + 366*uint32(numberOfLeapYears))

	// edge case: we are in january
	if daysPassedSinceJanFirst < daysInMonths[0] {
		return "Jan " + strconv.FormatUint(uint64(daysPassedSinceJanFirst+1), baseTen)
	}

	// check feb 29
	if yearIsLeap(currentYear) {
		daysInMonths[1]++
		if daysPassedSinceJanFirst >= daysInMonths[0]+daysInMonths[1] {
			daysPassedSinceJanFirst++
		}
	}

	var currentMonthIndex int = 0

	daysThisYear := daysInMonths

	for i := 0; i < 11; i++ {
		daysThisYear[i+1] += daysThisYear[i]
	}

	// get current month
	for i, val := range daysThisYear {
		if daysPassedSinceJanFirst >= val {
			currentMonthIndex = i
		}
	}

	// check if i'm in the first day of the next month
	if daysPassedSinceJanFirst >= daysThisYear[currentMonthIndex] {
		currentMonthIndex++
	}

	var currentDay uint32 = 1             // starting from january first
	currentDay += daysPassedSinceJanFirst // now I have the exact date of the year

	firstDayOfTheMonth := daysThisYear[currentMonthIndex-1] // the first day of the current month will be the value of the last month in this slice
	currentDay -= firstDayOfTheMonth                        // current day of the month = exact date of today - exact date of the first day of this month

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
