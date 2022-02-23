package chip8

func (emulator *Emulator) ClearScreen() {
	// TODO: clear the screen buffer
	emulator.PC += InstructionSize
}

func (emulator *Emulator) Return() {
	emulator.StackPop()
}

// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func (emulator *Emulator) Jump(nnn uint16) {
	emulator.PC = nnn
}

// Call subroutine at nnn.
//
// The interpreter increments the stack pointer, then puts the current PC
// on the top of the stack. The PC is then set to nnn.
func (emulator *Emulator) Call(nnn uint16) {
	emulator.StackPush()
	emulator.PC = nnn
}

// Skip next instruction if Vx = kk.
//
// The interpreter compares register Vx to kk, and if they are equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipEqual(x uint8, kk uint8) {
	if emulator.V[x] == kk {
		emulator.PC += InstructionSize * 2
	} else {
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipNotEqual(x uint8, kk uint8) {
	if emulator.V[x] != kk {
		emulator.PC += InstructionSize * 2
	} else {
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if Vx = Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipRegistersEqual(x uint8, y uint8) {
	if emulator.V[x] == emulator.V[y] {
		emulator.PC += InstructionSize * 2
	} else {
		emulator.PC += InstructionSize
	}
}

// Set Vx = kk.
//
// The interpreter puts the value kk into register Vx.
func (emulator *Emulator) LoadByte(x uint8, kk uint8) {
	emulator.V[x] = kk
}

// Set Vx = Vx + kk.
//
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func (emulator *Emulator) Add(x uint8, kk uint8) {
	emulator.V[x] += kk
}

// Set Vx = Vy.
//
// Stores the value of register Vy in register Vx.
func (emulator *Emulator) LoadRegister(x uint8, y uint8) {
	emulator.V[x] = emulator.V[y]
}
