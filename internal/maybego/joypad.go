package maybego

import "fmt"

type Joypad struct {
	prev_joypad byte
	directions  byte
	buttons     byte
}

type Direction byte
type Button byte

const JOYP uint16 = 0xFF00

const (
	DirectionRight Direction = iota
	DirectionLeft
	DirectionUp
	DirectionDown
)

const (
	ButtonA Button = iota
	ButtonB
	ButtonSelect
	ButtonStart
)

func NewJoypad() *Joypad {
	Write(0xFF00, 0x3F)
	return &Joypad{prev_joypad: Read(JOYP), directions: 0xF, buttons: 0xF}
}

func (joy *Joypad) updateControls() {
	new_joypad := Read(JOYP)
	// The buttons are only read on a write to JOYP
	// If there was no change, there is no need
	// changed := new_joypad != joy.prev_joypad

	// if !changed {
	// 	return
	// }

	directions_enabled := (new_joypad & 0x30) == 0x20
	buttons_enabled := (new_joypad & 0x30) == 0x10

	if !directions_enabled && !buttons_enabled {
		joy.prev_joypad = new_joypad | 0xF // all buttons disabled
		Write(JOYP, joy.prev_joypad)
		return
	}

	if directions_enabled {
		new_joypad |= joy.directions
		// joy.buttons = 0xF
		// fmt.Println("reset buttons")
	} else {
		new_joypad |= joy.buttons
		// joy.directions = 0xF
		// fmt.Println("reset directions")
	}

	bitmask := byte(0x1)
	for i := range 4 {
		if new_joypad&bitmask == 0 /*&& joy.prev_joypad&bitmask == 0*/ {
			RequestInterrupt(4)
		}
		bitmask <<= 1
	}

	joy.prev_joypad = new_joypad
	Write(JOYP, new_joypad)
}

func (joy *Joypad) setButton(b Button) {
	mask := byte(^(1 << b))
	joy.buttons &= mask
	fmt.Printf("Buttons: %X\n", joy.buttons)
}

func (joy *Joypad) resetButton(b Button) {
	mask := byte((1 << b))
	joy.buttons |= mask
	fmt.Printf("Buttons: %X\n", joy.buttons)
}

func (joy *Joypad) setDirection(d Direction) {
	mask := byte(^(1 << d))
	joy.directions &= mask
	fmt.Printf("Directions: %X\n", joy.directions)
}

func (joy *Joypad) resetDirection(d Direction) {
	mask := byte((1 << d))
	joy.directions |= mask
	fmt.Printf("Directions: %X\n", joy.directions)
}

