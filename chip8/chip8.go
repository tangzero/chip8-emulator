package chip8

type Emulator struct {
	V      [16]uint8   // general registers
	I      uint16      // address register
	SP     uint16      // stack pointer
	PC     uint16      // program counter
	Memory [4096]uint8 // 4KB of system RAM
	Stack  []uint8     // 32 bytes of stack. starting at 0x0EA0
	Screen []uint8     // 256 bytes of display buffer. starting at 0x0F00
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
	emulator.I = 0x00
	emulator.SP = 0x00
	emulator.PC = 0x0200
	emulator.Memory = [4096]uint8{}
	emulator.Stack = emulator.Memory[0x0EA0:]
	emulator.Screen = emulator.Memory[0x0F00:]
	emulator.Timer.Delay = 0
	emulator.Timer.Sound = 0
	copy(emulator.Memory[0x0200:], emulator.ROM)
}
