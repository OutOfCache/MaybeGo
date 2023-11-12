package maybego

import (
	"log"
	"os"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type logFlags struct {
	pc        bool
	registers bool
	flags     bool
}

type Logger struct {
	flags *logFlags
	debug bool
}

func NewLogger(debug bool, filename string) *Logger {
	logger := &Logger{flags: new(logFlags), debug: debug}
	log.SetFlags(0)
	if filename != "" {
		log.Printf("filename: %s", filename)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	}
	return logger
}

func (logger *Logger) SetPCFlag(flag bool) {
	logger.flags.pc = flag
}

func (logger *Logger) SetRegFlag(flag bool) {
	logger.flags.registers = flag
}

// bad name, I agree
func (logger *Logger) SetFlagsFlag(flag bool) {
	logger.flags.flags = flag
}

func (logger *Logger) LogRegisters(a byte, b byte, c byte, d byte, e byte, h byte, l byte) {
	if !logger.debug || !logger.flags.registers {
		return
	}

	log.Printf("%sA:%X BC:%X%X DE:%x%x HL:%x%x%s",
		Red, a, b, c, d, e, h, l, Reset)
}

func (logger *Logger) LogPC(pc uint16, cycles uint, op byte, arg0 byte, arg1 byte) {
	if !logger.debug || !logger.flags.pc {
		return
	}

	log.Printf("%sPC:%x (cy: %d) |0x%x: %x %x %x%s",
		Green, pc, cycles, pc, op, arg0, arg1, Reset)
}

func (logger *Logger) LogFlags(z bool, c bool, n bool, h bool, halt bool, ime bool) {
	if !logger.debug || !logger.flags.flags {
		return
	}
	zf := "-"
	if z {
		zf = "Z"
	}
	cf := "-"
	if c {
		cf = "C"
	}
	nf := "-"
	if n {
		nf = "N"
	}
	hf := "-"
	if h {
		hf = "H"
	}
	haltf := "-"
	if halt {
		haltf = "HALT"
	}
	imef := "-"
	if ime {
		imef = "IME"
	}

	log.Printf("%sF:%s%s%s%s %s %s%s",
		Cyan, zf, cf, nf, hf, haltf, imef, Reset)
}
