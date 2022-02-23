package chip8

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
