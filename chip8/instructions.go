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
	emulator.V[0xF] = (emulator.V[x] & 0b10000000) >> 7
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
	x = emulator.V[x]
	y = emulator.V[y]
	width := uint8(8)
	height := n

	emulator.V[0xF] = 0x00 // clean collision flag

	for yline := uint8(0); yline < height; yline++ {
		sprite := emulator.Memory[emulator.I+uint16(yline)]

		for xline := uint8(0); xline < width; xline++ {
			if (sprite & 0b10000000) != 0x00 {
				px, py := int(x+xline)%Width, int(y+yline)%Height
				index := emulator.Display.PixOffset(px, py) + 1 // color offset: 0:red, 1:green, 2:blue

				if emulator.Display.Pix[index] != 0x00 {
					emulator.V[0xF] = 0x01 // collision
				}

				emulator.Display.Pix[index] ^= 0xFF
			}
			sprite <<= 1
		}
	}
}

// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx
// is currently in the down position, PC is increased by 2.
func (emulator *Emulator) SkipKeyPressed(x uint8) {
	if emulator.KeyPressed(emulator.V[x]) {
		emulator.PC += InstructionSize
	}
}

// Skip next instruction if key with the value of Vx is not pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx
// is currently in the up position, PC is increased by 2.
func (emulator *Emulator) SkipKeyNotPressed(x uint8) {
	if !emulator.KeyPressed(emulator.V[x]) {
		emulator.PC += InstructionSize
	}
}

// Set Vx = delay timer value.
//
// The value of DT is placed into Vx.
func (emulator *Emulator) ReadDT(x uint8) {
	emulator.V[x] = emulator.DT
}

// Wait for a key press, store the value of the key in Vx.
//
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (emulator *Emulator) ReadKey(x uint8) {
	for {
		for key := uint8(0); key < 16; key++ {
			if emulator.KeyPressed(key) {
				emulator.V[x] = key
				return
			}
		}
	}
}

// Set delay timer = Vx.
//
// DT is set equal to the value of Vx.
func (emulator *Emulator) SetDT(x uint8) {
	emulator.DT = emulator.V[x]
}

// Set sound timer = Vx.
//
// ST is set equal to the value of Vx.
func (emulator *Emulator) SetST(x uint8) {
	emulator.ST = emulator.V[x]
}

// Set I = I + Vx.
//
// The values of I and Vx are added, and the results are stored in I.
func (emulator *Emulator) AddI(x uint8) {
	emulator.I += uint16(emulator.V[x])
}

// Set I = location of sprite for digit Vx.
//
// The value of I is set to the location for the hexadecimal sprite
// corresponding to the value of Vx.
func (emulator *Emulator) SetI(x uint8) {
	emulator.I = uint16(emulator.V[x]) * 5
}

// Store BCD representation of Vx in memory locations I, I+1, and I+2.
//
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory
// at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func (emulator *Emulator) LoadBCD(x uint8) {
	emulator.Memory[emulator.I] = emulator.V[x] / 100          // hundreds digit
	emulator.Memory[emulator.I+1] = (emulator.V[x] % 100) / 10 // tens digit
	emulator.Memory[emulator.I+2] = emulator.V[x] % 10         // ones digit
}

// Store registers V0 through Vx in memory starting at location I.
//
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
func (emulator *Emulator) StoreRegisters(x uint8) {
	copy(emulator.Memory[emulator.I:], emulator.V[:x+1])
}

// Read registers V0 through Vx from memory starting at location I.
//
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
func (emulator *Emulator) ReadRegisters(x uint8) {
	copy(emulator.V[:x+1], emulator.Memory[emulator.I:])
}
