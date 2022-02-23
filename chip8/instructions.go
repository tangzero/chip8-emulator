package chip8

func (emulator *Emulator) ClearScreen() {
	// TODO: clear the screen buffer
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
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipNotEqual(x uint8, kk uint8) {
	if emulator.V[x] != kk {
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if Vx = Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipRegistersEqual(x uint8, y uint8) {
	if emulator.V[x] == emulator.V[y] {
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
func (emulator *Emulator) AddByte(x uint8, kk uint8) {
	emulator.V[x] += kk
}

// Set Vx = Vy.
//
// Stores the value of register Vy in register Vx.
func (emulator *Emulator) LoadRegister(x uint8, y uint8) {
	emulator.V[x] = emulator.V[y]
}

// Set Vx = Vx OR Vy.
//
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
// A bitwise OR compares the corrseponding bits from two values, and if either bit is 1,
// then the same bit in the result is also 1. Otherwise, it is 0.
func (emulator *Emulator) Or(x uint8, y uint8) {
	emulator.V[x] |= emulator.V[y]
}

// Set Vx = Vx AND Vy.
//
// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
// A bitwise AND compares the corrseponding bits from two values, and if both bits are 1,
// then the same bit in the result is also 1. Otherwise, it is 0.
func (emulator *Emulator) And(x uint8, y uint8) {
	emulator.V[x] &= emulator.V[y]
}

// Set Vx = Vx XOR Vy.
//
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
// An exclusive OR compares the corrseponding bits from two values, and if the bits are not
// both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func (emulator *Emulator) Xor(x uint8, y uint8) {
	emulator.V[x] ^= emulator.V[y]
}

// Set Vx = Vx + Vy, set VF = carry.
//
// The values of Vx and Vy are added together.
// If the result is greater than 8 bits (i.e., > 255,) VF is set to 1,
// otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
func (emulator *Emulator) AddRegisters(x uint8, y uint8) {
	sum := uint16(emulator.V[x]) + uint16(emulator.V[y])
	emulator.V[x] = uint8(sum)
	emulator.V[0xF] = uint8(sum >> 8)
}

// SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
//
// If Vx > Vy, then VF is set to 1, otherwise 0.
// Then Vy is subtracted from Vx, and the results stored in Vx.
func (emulator *Emulator) Sub(x uint8, y uint8) {
	emulator.V[0xF] = map[bool]uint8{true: 1, false: 0}[emulator.V[x] > emulator.V[y]]
	emulator.V[x] -= emulator.V[y]
}

// Set Vx = Vx SHR 1.
//
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0.
// Then Vx is divided by 2.
func (emulator *Emulator) ShiftRight(x uint8) {
	emulator.V[0xF] = emulator.V[x] & 0b00000001
	emulator.V[x] >>= 1
}

// Set Vx = Vx SHL 1.
//
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0.
// Then Vx is multiplied by 2.
func (emulator *Emulator) ShiftLeft(x uint8) {
	emulator.V[0xF] = emulator.V[x] & 0b10000000
	emulator.V[x] <<= 1
}
