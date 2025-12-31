package maybego

var Memory [65536]byte

func InitMemory() {
	Memory[0xFF00] = 0xCF // init joypad input
}

func Read(adr uint16) byte {
	return Memory[adr]
}

func Write(adr uint16, val byte) {
	Memory[adr] = val
}
