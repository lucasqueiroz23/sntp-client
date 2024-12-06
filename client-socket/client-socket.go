package clientSocket

import (
	"bytes"
	"encoding/binary"
	"net"
	"sntp-client/error-handling"
)

type NtpPacket struct {
	Uli_vn_mode     uint8
	Ustratum        uint8
	Upoll           uint8
	Uprecision      uint8
	UrootDelay      uint32
	UrootDispersion uint32
	UrefId          uint32
	UrefTm_s        uint32
	UrefTm_f        uint32
	UorigTm_s       uint32
	UorigTm_f       uint32
	UrxTm_s         uint32
	UrxTm_f         uint32
	UtxTm_s         uint32
	TxTm_f          uint32
}

func buildPacketByteArray() []byte {
	packet := new(NtpPacket)
	packet.Uli_vn_mode = 0x1B

	message := bytes.NewBuffer(make([]byte, 0, 48))

	writeErr := binary.Write(message, binary.BigEndian, packet)
	if writeErr != nil {
		errorHandling.LogErrorAndExit(writeErr)
	}

	return message.Bytes()

}

func MakeRequest(IPAddress string) []byte {

	// the default port for the NTP protocol
	ntpPort := "123"
	fullAddress := IPAddress + ":" + ntpPort

	conn, dialErr := net.Dial("udp", fullAddress)

	if dialErr != nil {
		errorHandling.LogErrorAndExit(dialErr)
	}

	message := buildPacketByteArray()
	_, writeErr := conn.Write(message)

	if writeErr != nil {
		errorHandling.LogErrorAndExit(writeErr)
	}

	response := make([]byte, 48) // the size of a packet is 48 bytes
	_, readErr := conn.Read(response)
	if readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	conn.Close()

	return response
}
