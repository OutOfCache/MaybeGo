package main

import (
	"flag"
	"fmt"
	"github.com/outofcache/maybego/internal/maybego"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// var cpu *maybego.CPU
// var ppu *maybego.PPU
var emulator *maybego.Emulator
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

	emulator.LoadRom(&rom)
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

	// a := app.New()
	// display := a.NewWindow("Hello World")

	emulator = maybego.NewEmulator(logger)
	// ui = maybego.NewUI()
	// ppu.StartSDL()
	// ui.Setup()
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
		emulator.FetchDecodeExec()
		// ui.Update(emulator.GetPPU().GetCurrentFrame())

		// blarggs test
		// if maybego.Read(0xff02) == 0x81 {
		// 	c := maybego.Read(0xff01)
		// 	fmt.Printf("%c", c)
		// 	maybego.Write(0xff02, 0)

		// }
		// display.SetContent(widget.NewLabel("Hello World!"))
		// display.Show()
		// a.Run()
	}
	// ppu.EndSDL()
	// fmt.Print("quitting");
}
