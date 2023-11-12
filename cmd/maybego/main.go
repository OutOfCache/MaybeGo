package main

import (
	"flag"
	"fmt"
	"github.com/outofcache/maybego/internal/maybego"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var cpu *maybego.CPU

// main.go --debug --log-file=logs.txt --log=all

func loadROM() {
	if len(flag.Args()) != 1 {
		fmt.Println("Usage: go run main.go [-debug] [-logfile file] path/to/rom")
		os.Exit(1)
	}

	var path string = flag.Args()[0]

	rom, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File could not be read")
		fmt.Println(err)
		os.Exit(2)
	}

	for i, buffer := range rom {
		maybego.Write(uint16( /*0x100+*/ i), buffer)
	}
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

	cpu = maybego.NewCPU(logger)
	loadROM()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	quit := false

	go func() {
		<-sigs
		quit = true
	}()

	// FIXME: proper exit handling through SDL
	for !quit {
		cpu.Fetch()
		cycles := cpu.Decode()

		cpu.Handle_timer(cycles)

		// blarggs test
		if maybego.Read(0xff02) == 0x81 {
			c := maybego.Read(0xff01)
			fmt.Printf("%c", c)
			maybego.Write(0xff02, 0)
		}
	}
}
