package parser

import (
	"bytes"
	"encoding/binary"
	"sntp-client/client-socket"
	"sntp-client/error-handling"
)

func Parse(time []byte) *clientSocket.NtpPacket {
	response := new(clientSocket.NtpPacket)

	buf := bytes.NewReader(time)

	if readErr := binary.Read(buf, binary.BigEndian, response); readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	return response
}
