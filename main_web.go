//go:build js

package main

import "github.com/tangzero/chip8-emulator/chip8"

func LoadROM() chip8.ROM {
	data := DefaultROM
	name := "test_opcode"
	return chip8.ROM{Data: data, Name: name}
}
