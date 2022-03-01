package chip8_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tangzero/chip8-emulator/chip8"
)

func TestEmulator_Reset(t *testing.T) {
	emulator := chip8.NewEmulator(nil, nil)
	emulator.V[0x03] = 0xFF
	emulator.V[0x0F] = 0xBB

	emulator.Reset()

	assert.Equal(t, uint8(0x00), emulator.V[0x03])
	assert.Equal(t, uint8(0x00), emulator.V[0x0F])
}
