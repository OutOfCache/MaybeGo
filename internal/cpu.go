package maybego

type Registers struct {
	A  byte // can be combined to AF
	B  byte // BC, B hi
	C  byte // BC, C lo
	D  byte // DE, D hi
	E  byte // DE, E lo
	H  byte // HL, H hi
	L  byte // HL, L lo
	SP uint16
	PC uint16
}

type Flags struct {
	Z bool // zero flag
	C bool // carry flag
	N bool // sub flag
	H bool // half carry
}

var currentOpcode uint16;
var opcodes [256]func();

func cpu00() { // do I need parameters for args?
    return // NOP
}

func cpu01() int { // LD BC, u16
    Registers.B = Read(Registers.PC++)
    Registers.C = Read(Registers.PC++)
    Registers.PC++
}

func cpu02() int { // LD (BC), A
    address = (Registers.B << 8) + Registers.C
    Write(address, Registers.A)
    Registers.PC++
}
