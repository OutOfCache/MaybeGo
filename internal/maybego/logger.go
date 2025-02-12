package maybego

import (
	"log"
	"os"
	"io"
	"fmt"
	"bufio"
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
	writer io.Writer
}

func NewLogger(debug bool, filename string) *Logger {
	logger := &Logger{flags: new(logFlags), debug: debug, writer: os.Stdout}
	log.SetFlags(0)
	if filename != "" {
		log.Printf("filename: %s", filename)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}

		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""

		// log.SetOutput(file)
		logger.writer = bufio.NewWriter(file)

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

func (logger *Logger) LogRegisters(a byte, b byte, c byte, d byte, e byte, h byte, l byte, sp uint16) {
	if !logger.debug || !logger.flags.registers {
		return
	}

	// log.Printf("%sA:%02x BC:%02x%02x DE:%02x%02x HL:%02x%s",
	// 	Red, a, b, c, d, e, h, l, Reset)
	fmt.Fprintf(logger.writer, "%sA:%02x BC:%02X%02X DE:%02x%02x HL:%02x%02x SP:%4x%s ",
		Red, a, b, c, d, e, h, l, sp, Reset)
}

func (logger *Logger) LogPC(pc uint16, cycles uint, ppu byte, op byte, arg0 byte, arg1 byte) {
	if !logger.debug || !logger.flags.pc {
		return
	}

	// log.Printf("%sPC:%02x (cy: %d) ppu:+%d |0x%02x: %02x %02x %02x%s",
	// 	Green, pc, cycles*4, ppu, pc, op, arg0, arg1, Reset)
	fmt.Fprintf(logger.writer, "%sPC:%04x (cy: %d) ppu:+%d |0x%02x: %02x %02x %02x%s\n",
		Green, pc, cycles*4, ppu, pc, op, arg0, arg1, Reset)
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

	fmt.Fprintf(logger.writer, "%sF:%s%s%s%s %s %s%s ",
		Cyan, zf, nf, hf, cf, haltf, imef, Reset)
}

func (logger *Logger) LogValue(label string, val uint16) {
	if !logger.debug {
		return
	}

	fmt.Fprintf(logger.writer, "%s%s: %d%s\n",
		Green, label, val, Reset)
}
