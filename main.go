package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/tangzero/chip8-emulator/chip8"
)

//go:embed test_opcode.ch8
var DefaultROM []byte

const (
	ScreenScale = 20
	Width       = chip8.Width * ScreenScale
	Height      = chip8.Height * ScreenScale
)

type State = int

const (
	LoadingState State = iota
	RunningState
)

// Keyboard Layout
// 1 2 3 C
// 4 5 6 D
// 7 8 9 E
// A 0 B F
var KeyMapping = []ebiten.Key{
	ebiten.KeyX,
	ebiten.Key1,
	ebiten.Key2,
	ebiten.Key3,
	ebiten.KeyQ,
	ebiten.KeyW,
	ebiten.KeyE,
	ebiten.KeyA,
	ebiten.KeyS,
	ebiten.KeyD,
	ebiten.KeyZ,
	ebiten.KeyC,
	ebiten.Key4,
	ebiten.KeyR,
	ebiten.KeyF,
	ebiten.KeyV,
}

type GUI struct {
	State    State
	Emulator *chip8.Emulator
}

func (gui *GUI) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		gui.Emulator.Reset()
	}
	gui.Emulator.Update()
	return nil
}

func (gui *GUI) Draw(screen *ebiten.Image) {
	frame := ebiten.NewImageFromImage(gui.Emulator.Display)
	operation := new(ebiten.DrawImageOptions)
	operation.GeoM.Scale(ScreenScale, ScreenScale)
	screen.DrawImage(frame, operation)
}

func (gui *GUI) Layout(int, int) (int, int) {
	return Width, Height
}

func KeyPressed(key uint8) bool {
	return ebiten.IsKeyPressed(KeyMapping[key])
}

func SoundPlayer(sound []byte) (func(), func()) {
	stream, err := wav.DecodeWithSampleRate(chip8.SampleRate, bytes.NewReader(sound))
	assert(err)
	player, err := audio.NewContext(chip8.SampleRate).NewPlayer(stream)
	assert(err)
	player.SetVolume(0.5)
	return PlaySound(player), StopSound(player)
}

func PlaySound(player *audio.Player) func() {
	return func() {
		if player.IsPlaying() {
			return
		}
		player.Play()
	}
}

func StopSound(player *audio.Player) func() {
	return func() {
		player.Pause()
		assert(player.Rewind())
	}
}

func main() {
	rom := LoadROM()

	gui := GUI{}
	gui.State = LoadingState
	gui.Emulator = chip8.NewEmulator(KeyPressed, SoundPlayer)
	gui.Emulator.LoadROM(rom)

	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("CHIP-8 : " + rom.Name)

	assert(ebiten.RunGame(&gui))
}

func assert(err error) {
	if err != nil {
		log.Panic(err)
	}
}
