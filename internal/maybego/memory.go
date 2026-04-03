package maybego

import (
	"fmt"
	"os"
)

const (
	BANK     uint16 = 0xFF50 // Whether the boot rom has finished loading
	DMA      uint16 = 0xFF46 // OAM DMA
	OAMSTART uint16 = 0xFE00
	OAMEND   uint16 = 0xFE9F
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

func writeOamDma(src uint16) {
	// TODO: 160 cycles delay
	src_slice := Memory[(src << 8) : (src<<8)|0x9F]
	copy(Memory[0xFE00:0xFE9F], src_slice)
}

func Write(adr uint16, val byte) {
	if adr < 0x8000 && Memory[adr] != 0 {
		return
	}

	// OAM DMA: written val is the source address << 8
	// Dest: 0xFE00-0xFE9F
	if adr == DMA {
		writeOamDma(uint16(val))
	}
	Memory[adr] = val
}
