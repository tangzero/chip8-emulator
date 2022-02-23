package chip8

import "encoding/binary"

func (emulator *Emulator) StackPush() {
	if emulator.SP == StackAddress+StackSize*2 {
		panic("chip8: stack overflow")
	}
	binary.BigEndian.PutUint16(emulator.Memory[emulator.SP:], emulator.PC)
	emulator.SP += 2
}

func (emulator *Emulator) StackPop() {
	if emulator.SP == StackAddress {
		panic("chip8: nothing to pop from stack")
	}
	emulator.SP -= 2
	emulator.PC = binary.BigEndian.Uint16(emulator.Memory[emulator.SP:])
	binary.BigEndian.PutUint16(emulator.Memory[emulator.SP:], 0x00) // clean the stack position
}
