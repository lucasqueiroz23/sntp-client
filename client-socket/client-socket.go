package clientSocket

import (
	"encoding/json"
	"fmt"
	"net"
	"sntp-client/error-handling"
)

type ntpPacket struct {
	uli_vn_mode     int8
	ustratum        int8
	upoll           int8
	uprecision      int8
	urootDelay      int32
	urootDispersion int32
	urefId          int32
	urefTm_s        int32
	urefTm_f        int32
	uorigTm_s       int32
	uorigTm_f       int32
	urxTm_s         int32
	urxTm_f         int32
	utxTm_s         int32
	txTm_f          uint32
}

func buildPacketByteArray() []byte {
	packet := new(ntpPacket)
	packet.uli_vn_mode = 0x1B

	message, err := json.Marshal(packet)

	if err != nil {
		errorHandling.LogErrorAndExit(err)
	}

	return message
}

func MakeRequest(IPAddress string) {
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

	received := make([]byte, 1024)
	_, readErr := conn.Read(received)
	if readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	conn.Close()

	fmt.Println(received)
}
