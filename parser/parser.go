package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sntp-client/client-socket"
	"sntp-client/error-handling"
)

func Parse(serverResponse []byte) string {
	packet := getResponsePacket(serverResponse)

	fmt.Println(packet.TxTm_s)
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
