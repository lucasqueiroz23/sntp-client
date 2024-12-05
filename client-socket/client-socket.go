package clientSocket

import (
	"bytes"
	"encoding/binary"
	"net"
	"sntp-client/error-handling"
)

type ntpPacket struct {
	Uli_vn_mode     int8
	Ustratum        int8
	Upoll           int8
	Uprecision      int8
	UrootDelay      int32
	UrootDispersion int32
	UrefId          int32
	UrefTm_s        int32
	UrefTm_f        int32
	UorigTm_s       int32
	UorigTm_f       int32
	UrxTm_s         int32
	UrxTm_f         int32
	UtxTm_s         int32
	TxTm_f          uint32
}

func buildPacketByteArray() []byte {
	packet := new(ntpPacket)
	packet.Uli_vn_mode = 0x1B

	message := bytes.NewBuffer(make([]byte, 0, 48))

	writeErr := binary.Write(message, binary.BigEndian, packet)
	if writeErr != nil {
		errorHandling.LogErrorAndExit(writeErr)
	}

	return message.Bytes()

}

func MakeRequest(IPAddress string) []byte {

	message := buildPacketByteArray()

	// the default port for the NTP protocol
	ntpPort := "123"
	fullAddress := IPAddress + ":" + ntpPort

	conn, dialErr := net.Dial("udp", fullAddress)

	if dialErr != nil {
		errorHandling.LogErrorAndExit(dialErr)
	}

	_, writeErr := conn.Write(message)

	if writeErr != nil {
		errorHandling.LogErrorAndExit(writeErr)
	}

	response := make([]byte, 64)
	_, readErr := conn.Read(response)
	if readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	conn.Close()

	return response
}
