package main

import (
	// "../../internal/maybego"
	"fmt"
	"github.com/outofcache/maybego/internal/maybego"
	"io/ioutil"
	"os"
	"time"
)

var cpu *maybego.CPU

func loadROM() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go path/to/rom")
		os.Exit(1)
	}

	var path string = os.Args[1]

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
	cpu = maybego.NewCPU()
	cpuFreq := 1000.0 / 4194304 // 4.194304 MHz
	fmt.Print(cpuFreq)
	cpuCLK := time.NewTicker(time.Duration(cpuFreq) * time.Millisecond)

	loadROM()

	for _ = range cpuCLK.C {
		cpu.Fetch()
		cpu.Decode()

		// blarggs test
		if maybego.Read(0xff02) == 0x81 {
			c := maybego.Read(0xff01)
			fmt.Printf("%c", c)
			maybego.Write(0xff02, 0)
		}
	}
}
