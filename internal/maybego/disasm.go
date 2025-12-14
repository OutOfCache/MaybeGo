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
		func() string { return disasm.ldRegAddr("A", "BC") },    /* 0x0A */
		func() string { return disasm.dec("BC") },               /* 0x0B */
		func() string { return disasm.inc("C") },                /* 0x0C */
		func() string { return disasm.dec("C") },                /* 0x0D */
		func() string { return disasm.ldImm8("C") },             /* 0x0E */
		func() string { disasm.current_addr++; return "RRCA" },  /* 0x0F */
		func() string { disasm.current_addr++; return "STOP" },  /* 0x10 */
		func() string { return disasm.ldImm16("DE") },           /* 0x11 */
		func() string { return disasm.ldAddrReg("DE", "A") },    /* 0x12 */
		func() string { return disasm.inc("DE") },               /* 0x13 */
		func() string { return disasm.inc("D") },                /* 0x14 */
		func() string { return disasm.dec("D") },                /* 0x15 */
		func() string { return disasm.ldImm8("D") },             /* 0x16 */
		func() string { disasm.current_addr++; return "RLA" },   /* 0x17 */
		func() string { return disasm.jr() },                    /* 0x18 */
		func() string { return disasm.addReg("HL", "DE") },      /* 0x19 */
		func() string { return disasm.ldRegAddr("A", "DE") },    /* 0x1A */
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
		func() string { return disasm.ldRegAddr("A", "HL+") },   /* 0x2A */
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
		func() string { return disasm.ldRegAddr("A", "HL-") },   /* 0x3A */
		func() string { return disasm.dec("SP") },               /* 0x3B */
		func() string { return disasm.inc("A") },                /* 0x3C */
		func() string { return disasm.dec("A") },                /* 0x3D */
		func() string { return disasm.ldImm8("A") },             /* 0x3E */
		func() string { disasm.current_addr++; return "CCF" },   /* 0x3F */
		func() string { return disasm.ldRegReg("B", "B") },      /* 0x40 */
		func() string { return disasm.ldRegReg("B", "C") },      /* 0x41 */
		func() string { return disasm.ldRegReg("B", "D") },      /* 0x42 */
		func() string { return disasm.ldRegReg("B", "E") },      /* 0x43 */
		func() string { return disasm.ldRegReg("B", "H") },      /* 0x44 */
		func() string { return disasm.ldRegReg("B", "L") },      /* 0x45 */
		func() string { return disasm.ldRegAddr("B", "HL") },    /* 0x46 */
		func() string { return disasm.ldRegReg("B", "A") },      /* 0x47 */
		func() string { return disasm.ldRegReg("C", "B") },      /* 0x48 */
		func() string { return disasm.ldRegReg("C", "C") },      /* 0x49 */
		func() string { return disasm.ldRegReg("C", "D") },      /* 0x4A */
		func() string { return disasm.ldRegReg("C", "E") },      /* 0x4B */
		func() string { return disasm.ldRegReg("C", "H") },      /* 0x4C */
		func() string { return disasm.ldRegReg("C", "L") },      /* 0x4D */
		func() string { return disasm.ldRegAddr("C", "HL") },    /* 0x4E */
		func() string { return disasm.ldRegReg("C", "A") },      /* 0x4F */
		func() string { return disasm.ldRegReg("D", "B") },      /* 0x50 */
		func() string { return disasm.ldRegReg("D", "C") },      /* 0x51 */
		func() string { return disasm.ldRegReg("D", "D") },      /* 0x52 */
		func() string { return disasm.ldRegReg("D", "E") },      /* 0x53 */
		func() string { return disasm.ldRegReg("D", "H") },      /* 0x54 */
		func() string { return disasm.ldRegReg("D", "L") },      /* 0x55 */
		func() string { return disasm.ldRegAddr("D", "HL") },    /* 0x56 */
		func() string { return disasm.ldRegReg("D", "A") },      /* 0x57 */
		func() string { return disasm.ldRegReg("E", "B") },      /* 0x58 */
		func() string { return disasm.ldRegReg("E", "C") },      /* 0x59 */
		func() string { return disasm.ldRegReg("E", "D") },      /* 0x5A */
		func() string { return disasm.ldRegReg("E", "E") },      /* 0x5B */
		func() string { return disasm.ldRegReg("E", "H") },      /* 0x5C */
		func() string { return disasm.ldRegReg("E", "L") },      /* 0x5D */
		func() string { return disasm.ldRegAddr("E", "HL") },    /* 0x5E */
		func() string { return disasm.ldRegReg("E", "A") },      /* 0x5F */
		func() string { return disasm.ldRegReg("H", "B") },      /* 0x60 */
		func() string { return disasm.ldRegReg("H", "C") },      /* 0x61 */
		func() string { return disasm.ldRegReg("H", "D") },      /* 0x62 */
		func() string { return disasm.ldRegReg("H", "E") },      /* 0x63 */
		func() string { return disasm.ldRegReg("H", "H") },      /* 0x64 */
		func() string { return disasm.ldRegReg("H", "L") },      /* 0x65 */
		func() string { return disasm.ldRegAddr("H", "HL") },    /* 0x66 */
		func() string { return disasm.ldRegReg("H", "A") },      /* 0x67 */
		func() string { return disasm.ldRegReg("L", "B") },      /* 0x68 */
		func() string { return disasm.ldRegReg("L", "C") },      /* 0x69 */
		func() string { return disasm.ldRegReg("L", "D") },      /* 0x6A */
		func() string { return disasm.ldRegReg("L", "E") },      /* 0x6B */
		func() string { return disasm.ldRegReg("L", "H") },      /* 0x6C */
		func() string { return disasm.ldRegReg("L", "L") },      /* 0x6D */
		func() string { return disasm.ldRegAddr("L", "HL") },    /* 0x6E */
		func() string { return disasm.ldRegReg("L", "A") },      /* 0x6F */
		func() string { return disasm.ldAddrReg("HL", "B") },    /* 0x70 */
		func() string { return disasm.ldAddrReg("HL", "C") },    /* 0x71 */
		func() string { return disasm.ldAddrReg("HL", "D") },    /* 0x72 */
		func() string { return disasm.ldAddrReg("HL", "E") },    /* 0x73 */
		func() string { return disasm.ldAddrReg("HL", "H") },    /* 0x74 */
		func() string { return disasm.ldAddrReg("HL", "L") },    /* 0x75 */
		func() string { disasm.current_addr++; return "HALT" },  /* 0x76 */
		func() string { return disasm.ldAddrReg("HL", "A") },    /* 0x77 */
		func() string { return disasm.ldRegReg("A", "B") },      /* 0x78 */
		func() string { return disasm.ldRegReg("A", "C") },      /* 0x79 */
		func() string { return disasm.ldRegReg("A", "D") },      /* 0x7A */
		func() string { return disasm.ldRegReg("A", "E") },      /* 0x7B */
		func() string { return disasm.ldRegReg("A", "H") },      /* 0x7C */
		func() string { return disasm.ldRegReg("A", "L") },      /* 0x7D */
		func() string { return disasm.ldRegAddr("A", "HL") },    /* 0x7E */
		func() string { return disasm.ldRegReg("A", "A") },      /* 0x7F */
		func() string { return disasm.addReg("A", "B") },        /* 0x80 */
		func() string { return disasm.addReg("A", "C") },        /* 0x81 */
		func() string { return disasm.addReg("A", "D") },        /* 0x82 */
		func() string { return disasm.addReg("A", "E") },        /* 0x83 */
		func() string { return disasm.addReg("A", "H") },        /* 0x84 */
		func() string { return disasm.addReg("A", "L") },        /* 0x85 */
		func() string { return disasm.addReg("A", "(HL)") },     /* 0x86 */
		func() string { return disasm.addReg("A", "A") },        /* 0x87 */
		func() string { return disasm.adcReg("B") },             /* 0x88 */
		func() string { return disasm.adcReg("C") },             /* 0x89 */
		func() string { return disasm.adcReg("D") },             /* 0x8A */
		func() string { return disasm.adcReg("E") },             /* 0x8B */
		func() string { return disasm.adcReg("H") },             /* 0x8C */
		func() string { return disasm.adcReg("L") },             /* 0x8D */
		func() string { return disasm.adcReg("[HL]") },          /* 0x8E */
		func() string { return disasm.adcReg("A") },             /* 0x8F */
		func() string { return disasm.subReg("B") },             /* 0x90 */
		func() string { return disasm.subReg("C") },             /* 0x91 */
		func() string { return disasm.subReg("D") },             /* 0x92 */
		func() string { return disasm.subReg("E") },             /* 0x93 */
		func() string { return disasm.subReg("H") },             /* 0x94 */
		func() string { return disasm.subReg("L") },             /* 0x95 */
		func() string { return disasm.subReg("[HL]") },          /* 0x96 */
		func() string { return disasm.subReg("A") },             /* 0x97 */
		func() string { return disasm.sbcReg("B") },             /* 0x98 */
		func() string { return disasm.sbcReg("C") },             /* 0x99 */
		func() string { return disasm.sbcReg("D") },             /* 0x9A */
		func() string { return disasm.sbcReg("E") },             /* 0x9B */
		func() string { return disasm.sbcReg("H") },             /* 0x9C */
		func() string { return disasm.sbcReg("L") },             /* 0x9D */
		func() string { return disasm.sbcReg("[HL]") },          /* 0x9E */
		func() string { return disasm.sbcReg("A") },             /* 0x9F */
		func() string { return disasm.andReg("B") },             /* 0xA0 */
		func() string { return disasm.andReg("C") },             /* 0xA1 */
		func() string { return disasm.andReg("D") },             /* 0xA2 */
		func() string { return disasm.andReg("E") },             /* 0xA3 */
		func() string { return disasm.andReg("H") },             /* 0xA4 */
		func() string { return disasm.andReg("L") },             /* 0xA5 */
		func() string { return disasm.andReg("[HL]") },          /* 0xA6 */
		func() string { return disasm.andReg("A") },             /* 0xA7 */
		func() string { return disasm.xorReg("B") },             /* 0xA8 */
		func() string { return disasm.xorReg("C") },             /* 0xA9 */
		func() string { return disasm.xorReg("D") },             /* 0xAA */
		func() string { return disasm.xorReg("E") },             /* 0xAB */
		func() string { return disasm.xorReg("H") },             /* 0xAC */
		func() string { return disasm.xorReg("L") },             /* 0xAD */
		func() string { return disasm.xorReg("[HL]") },          /* 0xAE */
		func() string { return disasm.xorReg("A") },             /* 0xAF */
		func() string { return disasm.orReg("B") },              /* 0xB0 */
		func() string { return disasm.orReg("C") },              /* 0xB1 */
		func() string { return disasm.orReg("D") },              /* 0xB2 */
		func() string { return disasm.orReg("E") },              /* 0xB3 */
		func() string { return disasm.orReg("H") },              /* 0xB4 */
		func() string { return disasm.orReg("L") },              /* 0xB5 */
		func() string { return disasm.orReg("[HL]") },           /* 0xB6 */
		func() string { return disasm.orReg("A") },              /* 0xB7 */
		func() string { return disasm.cpReg("B") },              /* 0xB8 */
		func() string { return disasm.cpReg("C") },              /* 0xB9 */
		func() string { return disasm.cpReg("D") },              /* 0xBA */
		func() string { return disasm.cpReg("E") },              /* 0xBB */
		func() string { return disasm.cpReg("H") },              /* 0xBC */
		func() string { return disasm.cpReg("L") },              /* 0xBD */
		func() string { return disasm.cpReg("[HL]") },           /* 0xBE */
		func() string { return disasm.cpReg("A") },              /* 0xBF */
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
		func() string { disasm.current_addr++; return "DI" },    /* 0xF3 */
		func() string { return disasm.not_implemented() },       /* 0xF4 */
		func() string { return disasm.not_implemented() },       /* 0xF5 */
		func() string { return disasm.not_implemented() },       /* 0xF6 */
		func() string { return disasm.not_implemented() },       /* 0xF7 */
		func() string { return disasm.not_implemented() },       /* 0xF8 */
		func() string { return disasm.not_implemented() },       /* 0xF9 */
		func() string { return disasm.ldRegImm16("A") },         /* 0xFA */
		func() string { disasm.current_addr++; return "EI" },    /* 0xFB */
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
	return "LD [" + reg + "], " + reg + "\n"
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
	return "LD " + reg + ", [" + dis.printImm16At(dis.current_addr) + "]\n"
}

func (dis *Disasm) ldRegAddr(reg string, addr string) string {
	dis.current_addr++
	return "LD " + reg + ", [" + addr + "]\n"
}

func (dis *Disasm) ldRegReg(dst string, src string) string {
	dis.current_addr++
	return "LD " + dst + ", " + src + "\n"
}

func (dis *Disasm) jr() string {
	dis.current_addr++
	return "JR " + dis.printImm8At(dis.current_addr) + "\n"
}

func (dis *Disasm) jr_flag(flag string) string {
	dis.current_addr++
	return "JR " + flag + ", " + dis.printImm8At(dis.current_addr) + "\n"
}

func (dis *Disasm) addImm8(dst string) string {
	dis.current_addr++
	return "ADD " + dst + ", " + dis.printImm16At(dis.current_addr) + "\n"
}

func (dis *Disasm) addReg(dst string, src string) string {
	dis.current_addr++
	return "ADD " + dst + ", " + src + "\n"
}

func (dis *Disasm) adcReg(src string) string {
	dis.current_addr++
	return "ADC A, " + src + "\n"
}

func (dis *Disasm) subReg(src string) string {
	dis.current_addr++
	return "SUB A, " + src + "\n"
}

func (dis *Disasm) sbcReg(src string) string {
	dis.current_addr++
	return "SBC A, " + src + "\n"
}

func (dis *Disasm) andReg(src string) string {
	dis.current_addr++
	return "AND A, " + src + "\n"
}

func (dis *Disasm) xorReg(src string) string {
	dis.current_addr++
	return "XOR A, " + src + "\n"
}

func (dis *Disasm) orReg(src string) string {
	dis.current_addr++
	return "OR A, " + src + "\n"
}

func (dis *Disasm) cpReg(src string) string {
	dis.current_addr++
	return "CP A, " + src + "\n"
}

func (dis *Disasm) not_implemented() string {
	dis.current_addr++
	return "not implemented yet!\n"
}
