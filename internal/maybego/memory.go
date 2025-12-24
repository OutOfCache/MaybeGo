package maybego

var Memory [65536]byte

func Read(adr uint16) byte {
	return Memory[adr]
}

func Write(adr uint16, val byte) {
	if adr == 0xFF00 {
		val = 0xF
	}
	Memory[adr] = val
}
