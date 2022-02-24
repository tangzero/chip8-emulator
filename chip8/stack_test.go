package chip8_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tangzero/chip8-emulator/chip8"
)

func TestStack_Push(t *testing.T) {
	stack := chip8.NewStack()

	stack.Push(0xABCD)
	stack.Push(0xCDEF)
	stack.Push(0x9999)

	assert.Equal(t, []uint16{0xABCD, 0xCDEF, 0x9999}, stack.Values)
}

func TestStack_Push_Overflow(t *testing.T) {
	stack := chip8.NewStack()

	defer func() {
		assert.Equal(t, "chip8: stack overflow", recover())
	}()

	for {
		stack.Push(0xCAFE)
	}
}

func TestEmulator_StackPop(t *testing.T) {
	stack := chip8.NewStack()

	stack.Values = append(stack.Values, 0xCAFE)

	assert.Equal(t, uint16(0xCAFE), stack.Pop())
	assert.Equal(t, []uint16{}, stack.Values)
}

func TestEmulator_StackPop_Empty(t *testing.T) {
	stack := chip8.NewStack()

	defer func() {
		assert.Equal(t, "chip8: nothing to pop from stack", recover())
	}()

	stack.Pop()
}
