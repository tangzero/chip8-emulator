package chip8

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

// Clear the display.
func (emulator *Emulator) ClearScreen() {
	draw.Draw(emulator.Display, emulator.Display.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
}

// Return from a subroutine.
//
// The interpreter sets the program counter to the address at the top of the stack.
func (emulator *Emulator) Return() {
	emulator.PC = emulator.Stack.Pop()
}

// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func (emulator *Emulator) Jump(nnn uint16) {
	emulator.PC = nnn
}

// Call subroutine at nnn.
//
// Puts the current PC on the top of the stack.
// The PC is then set to nnn.
func (emulator *Emulator) Call(nnn uint16) {
	emulator.Stack.Push(emulator.PC)
	emulator.PC = nnn
}

// Skip next instruction if Vx = kk.
//
// The interpreter compares register Vx to kk, and if they are equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipEqualByte(x uint8, kk uint8) {
	if emulator.V[x] == kk {
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipNotEqualByte(x uint8, kk uint8) {
	if emulator.V[x] != kk {
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if Vx = Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (emulator *Emulator) SkipEqual(x uint8, y uint8) {
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
func (emulator *Emulator) Load(x uint8, y uint8) {
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
func (emulator *Emulator) Add(x uint8, y uint8) {
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
	emulator.V[0xF] = emulator.V[x] & 0b0000_0001
	emulator.V[x] >>= 1
}

// Set Vx = Vy - Vx, set VF = NOT borrow.
//
// If Vy > Vx, then VF is set to 1, otherwise 0.
// Then Vx is subtracted from Vy, and the results stored in Vx.
func (emulator *Emulator) SubN(x uint8, y uint8) {
	emulator.V[0xF] = map[bool]uint8{true: 1, false: 0}[emulator.V[y] > emulator.V[x]]
	emulator.V[x] = emulator.V[y] - emulator.V[x]
}

// Set Vx = Vx SHL 1.
//
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0.
// Then Vx is multiplied by 2.
func (emulator *Emulator) ShiftLeft(x uint8) {
	emulator.V[0xF] = emulator.V[x] & 0b1000_0000
	emulator.V[x] <<= 1
}

// Skip next instruction if Vx != Vy.
//
// The values of Vx and Vy are compared, and if they are not equal,
// the program counter is increased by 2.
func (emulator *Emulator) SkipNotEqual(x uint8, y uint8) {
	if emulator.V[x] != emulator.V[y] {
		emulator.PC += InstructionSize
	}
}

// Set I = nnn.
//
// The value of register I is set to nnn.
func (emulator *Emulator) LoadI(nnn uint16) {
	emulator.I = nnn
}

// Jump to location nnn + V0.
//
// The program counter is set to nnn plus the value of V0.
func (emulator *Emulator) JumpV0(nnn uint16) {
	emulator.PC = nnn + uint16(emulator.V[0])
}

// Set Vx = random byte AND kk.
//
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk.
// The results are stored in Vx. See instruction 8xy2 for more information on AND.
func (emulator *Emulator) Random(x uint8, kk uint8) {
	emulator.V[x] = uint8(rand.Intn(0xFF)) & kk
}

// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
//
// The interpreter reads n bytes from memory, starting at the address stored in I.
// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
// Sprites are XORed onto the existing screen. If this causes any pixels to be erased,
// VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it
// is outside the coordinates of the display, it wraps around to the opposite side of the screen.
func (emulator *Emulator) Draw(x uint8, y uint8, n uint8) {
	width := uint8(8)
	height := n

	emulator.V[0xF] = 0x00 // clean collision flag

	for row := uint8(0); row < height; row += 1 {
		sprite := emulator.Memory[emulator.I+uint16(row)]

		for col := uint8(0); col < width; col += 1 {
			if (sprite & 0b1000_0000) == 0x1 {
				px, py := int(emulator.V[x]+col), int(emulator.V[y]+row)

				// check for pixel collision
				if emulator.Display.At(px, py) != color.Black {
					emulator.V[0xF] = 0x01 // set collision flag
				}

				// draw pixel
				emulator.Display.Set(px, py, color.White)
			}

			// shift the sprite left 1. This will move the next next col/bit of the sprite into the first position.
			sprite <<= 1
		}
	}
}

// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx
// is currently in the down position, PC is increased by 2.
func (emulator *Emulator) SkipKeyPressed(x uint8) {
	// TODO: implement input
}

// Skip next instruction if key with the value of Vx is not pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx
// is currently in the up position, PC is increased by 2.
func (emulator *Emulator) SkipKeyNotPressed(x uint8) {
	// TODO: implement input
}
