package chip8_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tangzero/chip8-emulator/chip8"
)

func TestEmulator_StackPush(t *testing.T) {
	emulator := chip8.NewEmulator()

	emulator.PC = 0xABCD
	emulator.StackPush()

	assert.Equal(t, uint8(0xAB), emulator.Memory[chip8.StackAddress])
	assert.Equal(t, uint8(0xCD), emulator.Memory[chip8.StackAddress+1])
	assert.Equal(t, chip8.StackAddress+2, emulator.SP)
}

func TestEmulator_StackPop(t *testing.T) {
	emulator := chip8.NewEmulator()

	emulator.SP = chip8.StackAddress + 32
	emulator.Memory[chip8.StackAddress+30] = 0xEE
	emulator.Memory[chip8.StackAddress+31] = 0xFF
	emulator.StackPop()

	assert.Equal(t, uint16(0xEEFF), emulator.PC)
	assert.Equal(t, chip8.StackAddress+30, emulator.SP)
}
