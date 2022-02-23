package chip8

func (emulator *Emulator) ClearScreen() {
	// TODO: clear the screen buffer
	emulator.PC += InstructionSize
}

func (emulator *Emulator) Return() {
	emulator.StackPop()
}

func (emulator *Emulator) Jump(addr uint16) {
	emulator.PC = addr
}

func (emulator *Emulator) Call(addr uint16) {
	emulator.StackPush()
	emulator.PC = addr
}

func (emulator *Emulator) SkipEqual(x uint8, value uint8) {
	if emulator.V[x] == value {
		emulator.PC += InstructionSize * 2
	} else {
		emulator.PC += InstructionSize
	}
}

func (emulator *Emulator) SkipNotEqual(x uint8, value uint8) {
	if emulator.V[x] != value {
		emulator.PC += InstructionSize * 2
	} else {
		emulator.PC += InstructionSize
	}
}

func (emulator *Emulator) SkipRegistersEqual(x uint8, y uint8) {
	if emulator.V[x] == emulator.V[y] {
		emulator.PC += InstructionSize * 2
	} else {
		emulator.PC += InstructionSize
	}
}
