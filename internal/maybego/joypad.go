package maybego

const JOYP uint16 = 0xFF00

type Joypad struct {
	prev_joypad byte
}

func NewJoypad() *Joypad {
	return &Joypad{prev_joypad: Read(JOYP)}
}
