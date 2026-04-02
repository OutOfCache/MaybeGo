package maybego

import (
	"fmt"
	"os"
)

const (
	BANK uint16 = 0xFF50 // Whether the boot rom has finished loading
)

var Memory [65536]byte
var bootrom [0x100]byte

func InitMemory(boot *string) bool {
	Memory[JOYP] = 0xCF // init joypad input

	rom, err := os.ReadFile(*boot)
	if err != nil {
		fmt.Println("Bootrom could not be read. Skipping bootrom.")
		Memory[BANK] = 1
		return false
	}

	for i, buffer := range rom {
		bootrom[i] = buffer
	}

	return true
}

func Read(adr uint16) byte {
	if adr < 0x100 && Memory[BANK] == 0 {
		return bootrom[adr]
	}
	return Memory[adr]
}

func Write(adr uint16, val byte) {
	if adr < 0x8000 && Memory[adr] != 0 {
		return
	}
	Memory[adr] = val
}
