// package parser will get the contents of an NTP packet and
// parse it into a datetime in the format
// "Data/hora: [Day of the week] [Month] [Day of the month] [hh:mm:ss] [year]
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

// this must be printed before the datetime
const dateTimePrefix string = "Data/hora: "

// GetDate returns a datetime in the format specified in the head of this package.
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

// getResponsePacket reads the server response in bytes and
// converts it into a packet defined in the clientSocket package.
// If something wrong happens, it'll end the program.
func getResponsePacket(serverResponse []byte) *clientSocket.NtpPacket {
	response := new(clientSocket.NtpPacket)

	binaryResponse := bytes.NewReader(serverResponse)

	if readErr := binary.Read(binaryResponse, binary.BigEndian, response); readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	return response
}

// returns the current year based on the time passed in seconds since 1/1/1900 00:00:00
func getCurrentYear(timePassedInSeconds uint32) uint32 {
	return baseYear + uint32(float64(timePassedInSeconds)/(daysInAYear*hoursInADay*minutesInAnHour*secondsInAMinute))
}

// returns how many days have passed since 1/1/1900
func daysPastSince(timePassedInSeconds uint32) uint32 {
	return timePassedInSeconds / (uint32(hoursInADay) * uint32(minutesInAnHour) * uint32(secondsInAMinute))
}

// returns the day of the week given how many days have passed since 1/1/1900
func getDayOfTheWeek(daysSinceBaseDate uint32) string {
	return daysOfTheWeek[daysSinceBaseDate%daysInAWeek]
}

// returns the current time based on the time passed in seconds since 1/1/1900 00:00:00
func getCurrentTime(timePassedInSeconds uint32) string {

	seconds := timePassedInSeconds % uint32(secondsInAMinute)
	minutes := (timePassedInSeconds / uint32(secondsInAMinute)) % uint32(minutesInAnHour)
	hours := (timePassedInSeconds / (uint32(secondsInAMinute) * uint32(minutesInAnHour))) % uint32(hoursInADay)

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// gets the current day and month since 1/1/1900.
func getCurrentDay(daysSinceBaseDate uint32, currentYear uint32) string {
	daysInMonths := [12]uint32{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	yearsPassed := currentYear - baseYear
	numberOfLeapYears := countLeapYears(int(baseYear), int(currentYear))

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

	for i := 0; i < len(daysThisYear)-1; i++ {
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

// returns how many leap years have been since the year 0 up the year given as an argument
func leapYearsUpTo(year int) int {
	return (year / 4) - (year / 100) + (year / 400)
}

// returns how many leap years have in the given interval (start,end)
func countLeapYears(start, end int) int {
	return leapYearsUpTo(end) - leapYearsUpTo(start-1)
}

// returns true if a year is a leap year, false otherwise
func yearIsLeap(year uint32) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}
