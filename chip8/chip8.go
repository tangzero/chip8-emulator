package chip8

import (
	"encoding/binary"
	"image"
	"math"
)

const (
	MemorySize      = 4096 // 4KB of memory
	InstructionSize = 2
)

const (
	ProgramAddress = uint16(0x0200)
)

const (
	Width  = 64
	Height = 32
	FPS    = 60
)

type Emulator struct {
	V       [16]uint8         // general registers
	I       uint16            // address register
	PC      uint16            // program counter
	DT      uint8             // delay timer
	ST      uint8             // sound timer
	Stack   *Stack            // very simple stack
	Memory  [MemorySize]uint8 // 4KB of system RAM
	ROM     []uint8           // game rom
	Display *image.Gray       // display buffer
}

func NewEmulator() *Emulator {
	emulator := new(Emulator)
	emulator.Stack = NewStack()
	emulator.Display = image.NewGray(image.Rect(0, 0, Width, Height))
	emulator.Reset()
	return emulator
}

func (emulator *Emulator) Reset() {
	emulator.V = [16]uint8{}
	emulator.I = 0
	emulator.PC = ProgramAddress
	emulator.DT = 0
	emulator.ST = 0
	emulator.Memory = [MemorySize]uint8{}

	// clean stack
	emulator.Stack.Clear()

	// clean display
	emulator.ClearScreen()

	// reload ROM
	emulator.LoadROM(emulator.ROM)
}

func (emulator *Emulator) LoadROM(rom []uint8) {
	emulator.ROM = rom
	// load rom on memory
	copy(emulator.Memory[ProgramAddress:], emulator.ROM)
}

func (emulator *Emulator) Cycle() {
	pc := emulator.PC

	emulator.DT = uint8(math.Max(0, float64(emulator.DT)-1))
	emulator.ST = uint8(math.Max(0, float64(emulator.ST)-1))

	instruction := binary.BigEndian.Uint16(emulator.Memory[emulator.PC:])

	nnn := instruction & 0x0FFF
	n := uint8(instruction & 0x000F)
	kk := uint8(instruction & 0x00FF)
	x := uint8(instruction & 0x0F00 >> 8)
	y := uint8(instruction & 0x00F0 >> 4)

	switch instruction & 0xF000 >> 12 {
	case 0x0:
		switch instruction & 0x00FF {
		case 0xE0: // 00E0 - CLS
			emulator.ClearScreen()
		case 0xEE: // 00EE - RET
			emulator.Return()
		}
	case 0x1: // 1nnn - JP addr
		emulator.Jump(nnn)
	case 0x2: // 2nnn - CALL addr
		emulator.Call(nnn)
	case 0x3: // 3xkk - SE Vx, byte
		emulator.SkipEqualByte(x, kk)
	case 0x4: // 4xkk - SNE Vx, byte
		emulator.SkipNotEqualByte(x, kk)
	case 0x5: // 5xy0 - SE Vx, Vy
		emulator.SkipEqual(x, y)
	case 0x6: // 6xkk - LD Vx, byte
		emulator.LoadByte(x, kk)
	case 0x7: // 7xkk - ADD Vx, byte
		emulator.AddByte(x, kk)
	case 0x8:
		switch instruction & 0x000F {
		case 0x0: // 8xy0 - LD Vx, Vy
			emulator.Load(x, y)
		case 0x1: // 8xy1 - OR Vx, Vy
			emulator.Or(x, y)
		case 0x2: // 8xy2 - AND Vx, Vy
			emulator.And(x, y)
		case 0x3: // 8xy3 - XOR Vx, Vy
			emulator.Xor(x, y)
		case 0x4: // 8xy4 - ADD Vx, Vy
			emulator.Add(x, y)
		case 0x5: // 8xy5 - SUB Vx, Vy
			emulator.Sub(x, y)
		case 0x6: // 8xy6 - SHR Vx {, Vy}
			emulator.ShiftRight(x)
		case 0x7: // 8xy7 - SUBN Vx, Vy
			emulator.SubN(x, y)
		case 0xE: // 8xyE - SHL Vx {, Vy}
			emulator.ShiftLeft(x)
		}
	case 0x9: // 9xy0 - SNE Vx, Vy
		emulator.SkipNotEqual(x, y)
	case 0xA: // Annn - LD I, nnn
		emulator.LoadI(nnn)
	case 0xB: // Bnnn - JP V0, nnn
		emulator.JumpV0(nnn)
	case 0xC: // Cxkk - RND Vx, kk
		emulator.Random(x, kk)
	case 0xD: // Dxyn - DRW Vx, Vy, n
		emulator.Draw(x, y, n)
	case 0xE:
		switch instruction & 0x00FF {
		case 0x9E: // SKP Vx
			emulator.SkipKeyPressed(x)
		case 0xA1: // SKNP Vx
			emulator.SkipKeyNotPressed(x)
		}
	case 0xF:
		switch instruction & 0x00FF {
		case 0x07: // LD Vx, DT
			emulator.ReadDT(x)
		case 0x0A: // LD Vx, K
			emulator.ReadKey(x)
		case 0x15: // LD DT, Vx
			emulator.SetDT(x)
		case 0x18: // LD ST, Vx
			emulator.SetST(x)
		case 0x1E: // ADD I, Vx
			emulator.AddI(x)
		case 0x29: // LD F, Vx
			emulator.SetI(x)
		case 0x33: // LD B, Vx
			emulator.LoadBCD(x)
		case 0x55: // LD [I], Vx
			emulator.StoreRegisters(x)
		case 0x65: // LD Vx, [I]
			emulator.ReadRegisters(x)
		}
	}

	// if the program counter is unchanged and isn't a loop, read next instruction
	if emulator.PC == pc && nnn != pc {
		emulator.PC += InstructionSize
	}
}
