package clientSocket

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

func BuildPacket() *ntpPacket {
	packet := new(ntpPacket)
	packet.uli_vn_mode = 0x1B
	return packet
}
