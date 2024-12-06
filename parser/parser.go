package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sntp-client/client-socket"
	"sntp-client/error-handling"
)

const baseYear uint32 = 1900       // int NTP, the base date is jan 1, 1900
const daysInAYear float64 = 365.25 // the .25 accounts for years that have 366 days
const hoursInADay float64 = 24.0
const minutesInAnHour float64 = 60.0
const secondsInAMinute float64 = 60.0

func Parse(serverResponse []byte) string {
	packet := getResponsePacket(serverResponse)

	fmt.Println("ano atual: ", getCurrentYear(packet.TxTm_s))
	return ""
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
