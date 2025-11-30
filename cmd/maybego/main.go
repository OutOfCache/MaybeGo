package main

import (
	"flag"
	"fmt"
	"github.com/outofcache/maybego/internal/maybego"
	"os"
	"strings"
)

var ui *maybego.Interface

// main.go --debug --log-file=logs.txt --log=all

func loadROM() {
	if len(flag.Args()) != 1 {
		fmt.Println("Usage: go run main.go [-debug] [-logfile file] path/to/rom")
		os.Exit(1)
	}

	var path string = flag.Args()[0]

	rom, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("File could not be read")
		fmt.Println(err)
		os.Exit(2)
	}

	ui.LoadRom(&rom)
}

func main() {
	debugFlag := flag.Bool("debug", false, "enables logging")
	logFile := flag.String("logfile", "", "log output file")
	logContents := flag.String("logcontent", "", "what to log. Can be a combination of the following\npc\t\tlog pc and opcode information\nreg\t\tlog registers\nflags\tlog flags\nall\t\tlog everything")

	flag.Parse()
	logContentsSplit := strings.Split(*logContents, ",")

	logger := maybego.NewLogger(*debugFlag, *logFile)

	for _, c := range logContentsSplit {
		if c == "reg" || c == "all" {
			logger.SetRegFlag(true)
		}

		if c == "pc" || c == "all" {
			logger.SetPCFlag(true)
		}

		if c == "flags" || c == "all" {
			logger.SetFlagsFlag(true)
		}
	}

	ui = maybego.NewUI(logger)
	// TODO: optional argument
	loadROM()
	ui.Run()
}
