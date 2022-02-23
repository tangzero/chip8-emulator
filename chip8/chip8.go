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

func (emulator *Emulator) Step() {
	emulator.Timer.Delay = uint8(math.Max(0, float64(emulator.Timer.Delay)-1))
	emulator.Timer.Sound = uint8(math.Max(0, float64(emulator.Timer.Sound)-1))

	instruction := binary.BigEndian.Uint16(emulator.Memory[emulator.PC:])

	switch instruction & 0xF000 {
	case 0x1000: // 1nnn - JP addr
		emulator.Jump(instruction & 0x0FFF)
	case 0x2000: // 2nnn - CALL addr
		emulator.Call(instruction & 0x0FFF)
	case 0x3000: // 3xkk - SE Vx, byte
		emulator.SkipEqual(uint8(instruction&0x0F00>>8), uint8(instruction&0x00FF))
	case 0x4000: // 4xkk - SNE Vx, byte
		emulator.SkipNotEqual(uint8(instruction&0x0F00>>8), uint8(instruction&0x00FF))
	case 0x5000: // 5xy0 - SE Vx, Vy
		emulator.SkipRegistersEqual(uint8(instruction&0x0F00>>8), uint8(instruction&0x00F0>>4))
	}
}
