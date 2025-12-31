package maybego

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
	new_joypad := Read(JOYP) & 0xF0
	directions_enabled := (new_joypad & 0x30) == 0x20
	buttons_enabled := (new_joypad & 0x30) == 0x10

	if !directions_enabled && !buttons_enabled {
		joy.prev_joypad = new_joypad | 0xF // all buttons disabled
		Write(JOYP, joy.prev_joypad)
		return
	}

	if directions_enabled {
		new_joypad |= joy.directions
	} else {
		new_joypad |= joy.buttons
	}

	bitmask := byte(0x1)
	for range 4 {
		if new_joypad&bitmask == 0 {
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
}

func (joy *Joypad) resetButton(b Button) {
	mask := byte((1 << b))
	joy.buttons |= mask
}

func (joy *Joypad) setDirection(d Direction) {
	mask := byte(^(1 << d))
	joy.directions &= mask
}

func (joy *Joypad) resetDirection(d Direction) {
	mask := byte((1 << d))
	joy.directions |= mask
}
