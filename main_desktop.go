//go:build !js

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/tangzero/chip8-emulator/chip8"
)

func LoadROM() chip8.ROM {
	data := DefaultROM
	name := "test_opcode"

	if len(os.Args) > 1 {
		bytes, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		data = bytes
		name = strings.Split(path.Base(os.Args[1]), ".")[0]
	}

	return chip8.ROM{Data: data, Name: name}
}
