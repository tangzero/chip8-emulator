package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/tangzero/chip8-emulator/chip8"
)

const (
	ScreenScale = 20
	Width       = chip8.Width * ScreenScale
	Height      = chip8.Height * ScreenScale
)

type UI struct {
	Emulator *chip8.Emulator
}

func (ui *UI) Update(screen *ebiten.Image) error {
	ui.Emulator.Cycle()
	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	frame, err := ebiten.NewImageFromImage(ui.Emulator.Display, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Scale(ScreenScale, ScreenScale)

	screen.DrawImage(frame, op)
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func main() {
	rom, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ui := UI{Emulator: chip8.NewEmulator()}
	ui.Emulator.LoadROM(rom)

	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("CHIP 8")
	if err := ebiten.RunGame(&ui); err != nil {
		log.Fatal(err)
	}
}
