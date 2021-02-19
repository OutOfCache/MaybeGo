package maybego

import (
	"testing"
)

func TestWrite(t *testing.T) {
	var adr uint16 = 0x5623
	Write(adr, 0x08)

	if Memory[0x5623] != 0x08 {
		t.Error("Expected Memory[0x5623] to be 0x08")
	}
}

func TestRead(t *testing.T) {
	var adr uint16 = 0x5623
	Memory[adr] = 0x08
	if Read(adr) != 0x08 {
		t.Error("Expected Memory[0x5623] to be 0x08")
	}
}
