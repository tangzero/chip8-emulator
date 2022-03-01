package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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

type UI struct {
	Emulator     *chip8.Emulator
	State        State
	AudioContext *audio.Context
}

func (ui *UI) Run() {
	for {
		ui.Emulator.Cycle()
		time.Sleep(time.Millisecond * 2)
	}
}

func (ui *UI) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ui.Emulator.Reset()
	}

	switch ui.State {
	case LoadingState:
		go ui.Run()
		ui.State = RunningState
	case RunningState:
		ui.Emulator.UpdateTimers()
	}
	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	frame := ebiten.NewImageFromImage(ui.Emulator.Display)

	operation := new(ebiten.DrawImageOptions)
	operation.GeoM.Scale(ScreenScale, ScreenScale)

	screen.DrawImage(frame, operation)
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func (ui *UI) KeyPressed(key uint8) bool {
	return ebiten.IsKeyPressed(KeyMapping[key])
}

func (ui *UI) PlaySound(sound []byte) func() {
	player := ui.AudioContext.NewPlayerFromBytes(sound)
	player.SetVolume(0.3)
	player.Play()
	return func() { _ = player.Close() }
}

func main() {
	rom := LoadROM()

	ui := UI{}
	ui.Emulator = chip8.NewEmulator(ui.KeyPressed, ui.PlaySound)
	ui.State = LoadingState
	ui.AudioContext = audio.NewContext(44100)
	ui.Emulator.LoadROM(rom)

	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("CHIP-8 : " + rom.Name)

	if err := ebiten.RunGame(&ui); err != nil {
		log.Fatal(err)
	}
}
