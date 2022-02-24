package chip8

import (
	"encoding/binary"
	"math"
)

const (
	MemorySize      = 4096 // 4KB of memory
	StackSize       = 16
	InstructionSize = 2
)

const (
	ProgramAddress     = uint16(0x0200)
	StackAddress       = uint16(0x0EA0)
	VideoBufferAddress = uint16(0x0F00)
)

const (
	Width  = 64
	Height = 32
	FPS    = 60
)

type Emulator struct {
	V      [16]uint8         // general registers
	I      uint16            // address register
	SP     uint16            // stack pointer
	PC     uint16            // program counter
	Memory [MemorySize]uint8 // 4KB of system RAM
	Timer  struct {
		Delay uint8 // delay timer
		Sound uint8 // sound timer
	}
	ROM []uint8
}

func NewEmulator() *Emulator {
	emulator := new(Emulator)
	emulator.Reset()
	return emulator
}

func (emulator *Emulator) Reset() {
	emulator.V = [16]uint8{}
	emulator.I = 0
	emulator.SP = StackAddress
	emulator.PC = ProgramAddress
	emulator.Memory = [MemorySize]uint8{}
	emulator.Timer.Delay = 0
	emulator.Timer.Sound = 0
	copy(emulator.Memory[ProgramAddress:], emulator.ROM)
}

func (emulator *Emulator) Cycle() {
	pc := emulator.PC

	emulator.Timer.Delay = uint8(math.Max(0, float64(emulator.Timer.Delay)-1))
	emulator.Timer.Sound = uint8(math.Max(0, float64(emulator.Timer.Sound)-1))

	instruction := binary.BigEndian.Uint16(emulator.Memory[emulator.PC:])

	nnn := instruction & 0x0FFF
	// n := uint8(instruction & 0x000F)
	kk := uint8(instruction & 0x00FF)
	x := uint8(instruction & 0x0F00 >> 8)
	y := uint8(instruction & 0x00F0 >> 4)

	switch instruction & 0xF000 >> 12 {
	case 0x1: // 1nnn - JP addr
		emulator.Jump(nnn)
	case 0x2: // 2nnn - CALL addr
		emulator.Call(nnn)
	case 0x3: // 3xkk - SE Vx, byte
		emulator.SkipEqual(x, kk)
	case 0x4: // 4xkk - SNE Vx, byte
		emulator.SkipNotEqual(x, kk)
	case 0x5: // 5xy0 - SE Vx, Vy
		emulator.SkipRegistersEqual(x, y)
	case 0x6: // 6xkk - LD Vx, byte
		emulator.LoadByte(x, kk)
	case 0x7: // 7xkk - ADD Vx, byte
		emulator.AddByte(x, kk)
	case 0x8:
		switch instruction & 0x000F {
		case 0x0: // 8xy0 - LD Vx, Vy
			emulator.LoadRegister(x, y)
		case 0x1: // 8xy1 - OR Vx, Vy
			emulator.Or(x, y)
		case 0x2: // 8xy2 - AND Vx, Vy
			emulator.And(x, y)
		case 0x3: // 8xy3 - XOR Vx, Vy
			emulator.Xor(x, y)
		case 0x4: // 8xy4 - ADD Vx, Vy
			emulator.AddRegisters(x, y)
		case 0x5: // 8xy5 - SUB Vx, Vy
			emulator.Sub(x, y)
		case 0x6: // 8xy6 - SHR Vx {, Vy}
			emulator.ShiftRight(x)
		case 0x7: // 8xy7 - SUBN Vx, Vy
		case 0xE: // 8xyE - SHL Vx {, Vy}
			emulator.ShiftLeft(x)
		}
	case 0x9: // 9xy0 - SNE Vx, Vy
	case 0xA: // Annn - LD I, addr
	case 0xB: // Bnnn - JP V0, addr
	case 0xC: // Cxkk - RND Vx, byte
	case 0xD: // Dxyn - DRW Vx, Vy, nibble
	case 0xE:
		switch instruction & 0x00FF {
		case 0x9E: // SKP Vx
		case 0xA1: // SKNP Vx
		}
	case 0xF:
		switch instruction & 0x00FF {
		case 0x07: // LD Vx, DT
		case 0x0A: // LD Vx, K
		case 0x15: // LD DT, Vx
		case 0x18: // LD ST, Vx
		case 0x1E: // ADD I, Vx
		case 0x29: // LD F, Vx
		case 0x33: // LD B, Vx
		case 0x55: // LD [I], Vx
		case 0x65: // LD Vx, [I]
		}
	}

	// if the program counter is unchanged and isn't a loop, read next instruction
	if emulator.PC == pc && nnn != pc {
		emulator.PC += InstructionSize
	}
}
