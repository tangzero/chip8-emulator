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

type GUI struct {
	State    State
	Emulator *chip8.Emulator
}

func (gui *GUI) Run() {
	for {
		gui.Emulator.Cycle()
		time.Sleep(time.Millisecond * 2)
	}
}

func (gui *GUI) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		gui.Emulator.Reset()
	}
	switch gui.State {
	case LoadingState:
		go gui.Run()
		gui.State = RunningState
	case RunningState:
		gui.Emulator.UpdateTimers()
	}
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
	player := audio.NewContext(chip8.SampleRate).NewPlayerFromBytes(sound)
	player.SetVolume(0.3)
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
		player.Seek(0)
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

	if err := ebiten.RunGame(&gui); err != nil {
		log.Fatal(err)
	}
}
