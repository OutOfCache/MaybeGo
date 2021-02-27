package main

import (
	"../../internal/maybego"
	"fmt"
	"io/ioutil"
	"os"
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

	loadROM()

	for true {
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
