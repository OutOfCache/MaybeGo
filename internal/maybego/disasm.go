package maybego

import "fmt"

type Opcode struct {
	offset uint
	disasm string
}

type Disasm struct {
	file         *[]byte
	current_addr uint
	lines        []Opcode
	opcodes      [256]func() string
	cbOps        [256]func() string
}

func NewDisasm() *Disasm {
	disasm := &Disasm{current_addr: 0x0}
	disasm.opcodes = [256]func() string{
		func() string { disasm.current_addr++; return "nop\n" }, /* 0x00 */
		func() string { return disasm.ldImm16("BC") },           /* 0x01 */
		func() string { return disasm.ldAddrReg("BC", "A") },    /* 0x02 */
		func() string { return disasm.inc("BC") },               /* 0x03 */
		func() string { return disasm.inc("B") },                /* 0x04 */
		func() string { return disasm.dec("B") },                /* 0x05 */
		func() string { return disasm.ldImm8("B") },             /* 0x06 */
		func() string { disasm.current_addr++; return "RLCA" },  /* 0x07 */
		func() string { return disasm.ldImm16Reg("SP") },        /* 0x08 */
		func() string { return disasm.addReg("HL", "BC") },      /* 0x09 */
		func() string { return disasm.not_implemented() },       /* 0x0A */
		func() string { return disasm.dec("BC") },               /* 0x0B */
		func() string { return disasm.inc("C") },                /* 0x0C */
		func() string { return disasm.dec("C") },                /* 0x0D */
		func() string { return disasm.ldImm8("C") },             /* 0x0E */
		func() string { disasm.current_addr++; return "RRCA" },  /* 0x0F */
		func() string { return disasm.not_implemented() },       /* 0x10 */
		func() string { return disasm.ldImm16("DE") },           /* 0x11 */
		func() string { return disasm.ldAddrReg("DE", "A") },    /* 0x12 */
		func() string { return disasm.inc("DE") },               /* 0x13 */
		func() string { return disasm.inc("D") },                /* 0x14 */
		func() string { return disasm.dec("D") },                /* 0x15 */
		func() string { return disasm.ldImm8("D") },             /* 0x16 */
		func() string { disasm.current_addr++; return "RLA" },   /* 0x17 */
		func() string { return disasm.jr() },                    /* 0x18 */
		func() string { return disasm.addReg("HL", "DE") },      /* 0x19 */
		func() string { return disasm.not_implemented() },       /* 0x1A */
		func() string { return disasm.dec("DE") },               /* 0x1B */
		func() string { return disasm.inc("E") },                /* 0x1C */
		func() string { return disasm.dec("E") },                /* 0x1D */
		func() string { return disasm.ldImm8("E") },             /* 0x1E */
		func() string { disasm.current_addr++; return "RRA" },   /* 0x1F */
		func() string { return disasm.jr_flag("NZ") },           /* 0x20 */
		func() string { return disasm.ldImm16("HL") },           /* 0x21 */
		func() string { return disasm.ldAddrReg("HL+", "A") },   /* 0x22 */
		func() string { return disasm.inc("HL") },               /* 0x23 */
		func() string { return disasm.inc("H") },                /* 0x24 */
		func() string { return disasm.dec("H") },                /* 0x25 */
		func() string { return disasm.ldImm8("H") },             /* 0x26 */
		func() string { disasm.current_addr++; return "DAA" },   /* 0x27 */
		func() string { return disasm.jr_flag("Z") },            /* 0x28 */
		func() string { return disasm.addReg("HL", "HL") },      /* 0x29 */
		func() string { return disasm.not_implemented() },       /* 0x2A */
		func() string { return disasm.dec("HL") },               /* 0x2B */
		func() string { return disasm.inc("L") },                /* 0x2C */
		func() string { return disasm.dec("L") },                /* 0x2D */
		func() string { return disasm.ldImm8("L") },             /* 0x2E */
		func() string { disasm.current_addr++; return "CPL" },   /* 0x2F */
		func() string { return disasm.jr_flag("NC") },           /* 0x30 */
		func() string { return disasm.ldImm16("SP") },           /* 0x31 */
		func() string { return disasm.ldAddrReg("HL-", "A") },   /* 0x32 */
		func() string { return disasm.inc("SP") },               /* 0x33 */
		func() string { return disasm.inc("(HL)") },             /* 0x34 */
		func() string { return disasm.dec("(HL)") },             /* 0x35 */
		func() string { return disasm.ldImm8("(HL)") },          /* 0x36 */
		func() string { disasm.current_addr++; return "SCF" },   /* 0x37 */
		func() string { return disasm.jr_flag("C") },            /* 0x38 */
		func() string { return disasm.addReg("HL", "SP") },      /* 0x39 */
		func() string { return disasm.not_implemented() },       /* 0x3A */
		func() string { return disasm.dec("SP") },               /* 0x3B */
		func() string { return disasm.inc("A") },                /* 0x3C */
		func() string { return disasm.dec("A") },                /* 0x3D */
		func() string { return disasm.ldImm8("A") },             /* 0x3E */
		func() string { disasm.current_addr++; return "CCF" },   /* 0x3F */
		func() string { return disasm.not_implemented() },       /* 0x40 */
		func() string { return disasm.not_implemented() },       /* 0x41 */
		func() string { return disasm.not_implemented() },       /* 0x42 */
		func() string { return disasm.not_implemented() },       /* 0x43 */
		func() string { return disasm.not_implemented() },       /* 0x44 */
		func() string { return disasm.not_implemented() },       /* 0x45 */
		func() string { return disasm.not_implemented() },       /* 0x46 */
		func() string { return disasm.not_implemented() },       /* 0x47 */
		func() string { return disasm.not_implemented() },       /* 0x48 */
		func() string { return disasm.not_implemented() },       /* 0x49 */
		func() string { return disasm.not_implemented() },       /* 0x4A */
		func() string { return disasm.not_implemented() },       /* 0x4B */
		func() string { return disasm.not_implemented() },       /* 0x4C */
		func() string { return disasm.not_implemented() },       /* 0x4D */
		func() string { return disasm.not_implemented() },       /* 0x4E */
		func() string { return disasm.not_implemented() },       /* 0x4F */
		func() string { return disasm.not_implemented() },       /* 0x50 */
		func() string { return disasm.not_implemented() },       /* 0x51 */
		func() string { return disasm.not_implemented() },       /* 0x52 */
		func() string { return disasm.not_implemented() },       /* 0x53 */
		func() string { return disasm.not_implemented() },       /* 0x54 */
		func() string { return disasm.not_implemented() },       /* 0x55 */
		func() string { return disasm.not_implemented() },       /* 0x56 */
		func() string { return disasm.not_implemented() },       /* 0x57 */
		func() string { return disasm.not_implemented() },       /* 0x58 */
		func() string { return disasm.not_implemented() },       /* 0x59 */
		func() string { return disasm.not_implemented() },       /* 0x5A */
		func() string { return disasm.not_implemented() },       /* 0x5B */
		func() string { return disasm.not_implemented() },       /* 0x5C */
		func() string { return disasm.not_implemented() },       /* 0x5D */
		func() string { return disasm.not_implemented() },       /* 0x5E */
		func() string { return disasm.not_implemented() },       /* 0x5F */
		func() string { return disasm.not_implemented() },       /* 0x60 */
		func() string { return disasm.not_implemented() },       /* 0x61 */
		func() string { return disasm.not_implemented() },       /* 0x62 */
		func() string { return disasm.not_implemented() },       /* 0x63 */
		func() string { return disasm.not_implemented() },       /* 0x64 */
		func() string { return disasm.not_implemented() },       /* 0x65 */
		func() string { return disasm.not_implemented() },       /* 0x66 */
		func() string { return disasm.not_implemented() },       /* 0x67 */
		func() string { return disasm.not_implemented() },       /* 0x68 */
		func() string { return disasm.not_implemented() },       /* 0x69 */
		func() string { return disasm.not_implemented() },       /* 0x6A */
		func() string { return disasm.not_implemented() },       /* 0x6B */
		func() string { return disasm.not_implemented() },       /* 0x6C */
		func() string { return disasm.not_implemented() },       /* 0x6D */
		func() string { return disasm.not_implemented() },       /* 0x6E */
		func() string { return disasm.not_implemented() },       /* 0x6F */
		func() string { return disasm.not_implemented() },       /* 0x70 */
		func() string { return disasm.not_implemented() },       /* 0x71 */
		func() string { return disasm.not_implemented() },       /* 0x72 */
		func() string { return disasm.not_implemented() },       /* 0x73 */
		func() string { return disasm.not_implemented() },       /* 0x74 */
		func() string { return disasm.not_implemented() },       /* 0x75 */
		func() string { return disasm.not_implemented() },       /* 0x76 */
		func() string { return disasm.not_implemented() },       /* 0x77 */
		func() string { return disasm.not_implemented() },       /* 0x78 */
		func() string { return disasm.not_implemented() },       /* 0x79 */
		func() string { return disasm.not_implemented() },       /* 0x7A */
		func() string { return disasm.not_implemented() },       /* 0x7B */
		func() string { return disasm.not_implemented() },       /* 0x7C */
		func() string { return disasm.not_implemented() },       /* 0x7D */
		func() string { return disasm.not_implemented() },       /* 0x7E */
		func() string { return disasm.not_implemented() },       /* 0x7F */
		func() string { return disasm.addReg("A", "B") },        /* 0x80 */
		func() string { return disasm.addReg("A", "C") },        /* 0x81 */
		func() string { return disasm.addReg("A", "D") },        /* 0x82 */
		func() string { return disasm.addReg("A", "E") },        /* 0x83 */
		func() string { return disasm.addReg("A", "H") },        /* 0x84 */
		func() string { return disasm.addReg("A", "L") },        /* 0x85 */
		func() string { return disasm.addReg("A", "(HL)") },     /* 0x86 */
		func() string { return disasm.addReg("A", "A") },        /* 0x87 */
		func() string { return disasm.not_implemented() },       /* 0x88 */
		func() string { return disasm.not_implemented() },       /* 0x89 */
		func() string { return disasm.not_implemented() },       /* 0x8A */
		func() string { return disasm.not_implemented() },       /* 0x8B */
		func() string { return disasm.not_implemented() },       /* 0x8C */
		func() string { return disasm.not_implemented() },       /* 0x8D */
		func() string { return disasm.not_implemented() },       /* 0x8E */
		func() string { return disasm.not_implemented() },       /* 0x8F */
		func() string { return disasm.not_implemented() },       /* 0x90 */
		func() string { return disasm.not_implemented() },       /* 0x91 */
		func() string { return disasm.not_implemented() },       /* 0x92 */
		func() string { return disasm.not_implemented() },       /* 0x93 */
		func() string { return disasm.not_implemented() },       /* 0x94 */
		func() string { return disasm.not_implemented() },       /* 0x95 */
		func() string { return disasm.not_implemented() },       /* 0x96 */
		func() string { return disasm.not_implemented() },       /* 0x97 */
		func() string { return disasm.not_implemented() },       /* 0x98 */
		func() string { return disasm.not_implemented() },       /* 0x99 */
		func() string { return disasm.not_implemented() },       /* 0x9A */
		func() string { return disasm.not_implemented() },       /* 0x9B */
		func() string { return disasm.not_implemented() },       /* 0x9C */
		func() string { return disasm.not_implemented() },       /* 0x9D */
		func() string { return disasm.not_implemented() },       /* 0x9E */
		func() string { return disasm.not_implemented() },       /* 0x9F */
		func() string { return disasm.not_implemented() },       /* 0xA0 */
		func() string { return disasm.not_implemented() },       /* 0xA1 */
		func() string { return disasm.not_implemented() },       /* 0xA2 */
		func() string { return disasm.not_implemented() },       /* 0xA3 */
		func() string { return disasm.not_implemented() },       /* 0xA4 */
		func() string { return disasm.not_implemented() },       /* 0xA5 */
		func() string { return disasm.not_implemented() },       /* 0xA6 */
		func() string { return disasm.not_implemented() },       /* 0xA7 */
		func() string { return disasm.not_implemented() },       /* 0xA8 */
		func() string { return disasm.not_implemented() },       /* 0xA9 */
		func() string { return disasm.not_implemented() },       /* 0xAA */
		func() string { return disasm.not_implemented() },       /* 0xAB */
		func() string { return disasm.not_implemented() },       /* 0xAC */
		func() string { return disasm.not_implemented() },       /* 0xAD */
		func() string { return disasm.not_implemented() },       /* 0xAE */
		func() string { return disasm.not_implemented() },       /* 0xAF */
		func() string { return disasm.not_implemented() },       /* 0xB0 */
		func() string { return disasm.not_implemented() },       /* 0xB1 */
		func() string { return disasm.not_implemented() },       /* 0xB2 */
		func() string { return disasm.not_implemented() },       /* 0xB3 */
		func() string { return disasm.not_implemented() },       /* 0xB4 */
		func() string { return disasm.not_implemented() },       /* 0xB5 */
		func() string { return disasm.not_implemented() },       /* 0xB6 */
		func() string { return disasm.not_implemented() },       /* 0xB7 */
		func() string { return disasm.not_implemented() },       /* 0xB8 */
		func() string { return disasm.not_implemented() },       /* 0xB9 */
		func() string { return disasm.not_implemented() },       /* 0xBA */
		func() string { return disasm.not_implemented() },       /* 0xBB */
		func() string { return disasm.not_implemented() },       /* 0xBC */
		func() string { return disasm.not_implemented() },       /* 0xBD */
		func() string { return disasm.not_implemented() },       /* 0xBE */
		func() string { return disasm.not_implemented() },       /* 0xBF */
		func() string { return disasm.not_implemented() },       /* 0xC0 */
		func() string { return disasm.not_implemented() },       /* 0xC1 */
		func() string { return disasm.not_implemented() },       /* 0xC2 */
		func() string { return disasm.not_implemented() },       /* 0xC3 */
		func() string { return disasm.not_implemented() },       /* 0xC4 */
		func() string { return disasm.not_implemented() },       /* 0xC5 */
		func() string { return disasm.addImm8("A") },            /* 0xC6 */
		func() string { return disasm.not_implemented() },       /* 0xC7 */
		func() string { return disasm.not_implemented() },       /* 0xC8 */
		func() string { return disasm.not_implemented() },       /* 0xC9 */
		func() string { return disasm.not_implemented() },       /* 0xCA */
		func() string { return disasm.not_implemented() },       /* 0xCB */
		func() string { return disasm.not_implemented() },       /* 0xCC */
		func() string { return disasm.not_implemented() },       /* 0xCD */
		func() string { return disasm.not_implemented() },       /* 0xCE */
		func() string { return disasm.not_implemented() },       /* 0xCF */
		func() string { return disasm.not_implemented() },       /* 0xD0 */
		func() string { return disasm.not_implemented() },       /* 0xD1 */
		func() string { return disasm.not_implemented() },       /* 0xD2 */
		func() string { return disasm.not_implemented() },       /* 0xD3 */
		func() string { return disasm.not_implemented() },       /* 0xD4 */
		func() string { return disasm.not_implemented() },       /* 0xD5 */
		func() string { return disasm.not_implemented() },       /* 0xD6 */
		func() string { return disasm.not_implemented() },       /* 0xD7 */
		func() string { return disasm.not_implemented() },       /* 0xD8 */
		func() string { return disasm.not_implemented() },       /* 0xD9 */
		func() string { return disasm.not_implemented() },       /* 0xDA */
		func() string { return disasm.not_implemented() },       /* 0xDB */
		func() string { return disasm.not_implemented() },       /* 0xDC */
		func() string { return disasm.not_implemented() },       /* 0xDD */
		func() string { return disasm.not_implemented() },       /* 0xDE */
		func() string { return disasm.not_implemented() },       /* 0xDF */
		func() string { return disasm.not_implemented() },       /* 0xE0 */
		func() string { return disasm.not_implemented() },       /* 0xE1 */
		func() string { return disasm.not_implemented() },       /* 0xE2 */
		func() string { return disasm.not_implemented() },       /* 0xE3 */
		func() string { return disasm.not_implemented() },       /* 0xE4 */
		func() string { return disasm.not_implemented() },       /* 0xE5 */
		func() string { return disasm.not_implemented() },       /* 0xE6 */
		func() string { return disasm.not_implemented() },       /* 0xE7 */
		func() string { return disasm.addImm8("SP") },           /* 0xE8 */
		func() string { return disasm.not_implemented() },       /* 0xE9 */
		func() string { return disasm.ldImm16Reg("A") },         /* 0xEA */
		func() string { return disasm.not_implemented() },       /* 0xEB */
		func() string { return disasm.not_implemented() },       /* 0xEC */
		func() string { return disasm.not_implemented() },       /* 0xED */
		func() string { return disasm.not_implemented() },       /* 0xEE */
		func() string { return disasm.not_implemented() },       /* 0xEF */
		func() string { return disasm.not_implemented() },       /* 0xF0 */
		func() string { return disasm.not_implemented() },       /* 0xF1 */
		func() string { return disasm.not_implemented() },       /* 0xF2 */
		func() string { return disasm.not_implemented() },       /* 0xF3 */
		func() string { return disasm.not_implemented() },       /* 0xF4 */
		func() string { return disasm.not_implemented() },       /* 0xF5 */
		func() string { return disasm.not_implemented() },       /* 0xF6 */
		func() string { return disasm.not_implemented() },       /* 0xF7 */
		func() string { return disasm.not_implemented() },       /* 0xF8 */
		func() string { return disasm.not_implemented() },       /* 0xF9 */
		func() string { return disasm.ldRegImm16("A") },         /* 0xFA */
		func() string { return disasm.not_implemented() },       /* 0xFB */
		func() string { return disasm.not_implemented() },       /* 0xFC */
		func() string { return disasm.not_implemented() },       /* 0xFD */
		func() string { return disasm.not_implemented() },       /* 0xFE */
		func() string { return disasm.not_implemented() },       /* 0xFF */
	}

	return disasm
}

func (dis *Disasm) SetFile(file *[]byte) {
	dis.file = file
}

func (dis *Disasm) Disassemble() {
	for dis.current_addr < uint(len(*dis.file)) {
		opc := (*dis.file)[dis.current_addr]
		fmt.Printf("0x%X\t| %s", dis.current_addr, dis.opcodes[opc]())
	}
}

func (dis *Disasm) printImm8At(addr uint) string {
	dis.current_addr++
	return fmt.Sprintf("%02X %02X", (*dis.file)[addr])
}

func (dis *Disasm) printImm16At(start_addr uint) string {
	dis.current_addr += 2
	return fmt.Sprintf("%02X %02X", (*dis.file)[start_addr], (*dis.file)[start_addr+1])
}

func (dis *Disasm) ldImm16(reg string) string {
	dis.current_addr++
	return "LD " + reg + " " + dis.printImm16At(dis.current_addr) + "\n"
}

func (dis *Disasm) ldImm8(reg string) string {
	dis.current_addr++
	return "LD " + reg + " " + dis.printImm8At(dis.current_addr) + "\n"
}

func (dis *Disasm) ldAddrReg(addr string, reg string) string {
	dis.current_addr++
	return "LD (" + reg + ") " + reg + "\n"
}

func (dis *Disasm) inc(reg string) string {
	dis.current_addr++
	return "INC " + reg + "\n"
}

func (dis *Disasm) dec(reg string) string {
	dis.current_addr++
	return "DEC " + reg + "\n"
}

// Load [a16], reg
func (dis *Disasm) ldImm16Reg(reg string) string {
	dis.current_addr++
	return "LD [" + dis.printImm16At(dis.current_addr) + "], " + reg + "\n"
}

// Load reg, [a16]
func (dis *Disasm) ldRegImm16(reg string) string {
	dis.current_addr++
	return "LD " + reg + ", [" + dis.printImm16At(dis.current_addr) + "] " + "\n"
}

func (dis *Disasm) jr() string {
	dis.current_addr++
	return "JR " + dis.printImm8At(dis.current_addr) + "\n"
}

func (dis *Disasm) jr_flag(flag string) string {
	dis.current_addr++
	return "JR " + flag + ", " + dis.printImm8At(dis.current_addr) + "\n"
}

func (dis *Disasm) addReg(dst string, src string) string {
	dis.current_addr++
	return "ADD " + dst + ", " + src + "\n"
}

func (dis *Disasm) addImm8(dst string) string {
	dis.current_addr++
	return "ADD " + dst + ", " + dis.printImm16At(dis.current_addr) + "\n"
}

func (dis *Disasm) not_implemented() string {
	dis.current_addr++
	return "not implemented yet!\n"
}
