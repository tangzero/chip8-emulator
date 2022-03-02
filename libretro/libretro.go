package main

/*
#include "libretro.h"

typedef struct retro_system_info retro_system_info;
typedef struct retro_system_av_info retro_system_av_info;
typedef struct retro_game_info retro_game_info;

void VideoRefresh(const void *data, unsigned width, unsigned height, size_t pitch);
void InputPoll(void);
int16_t InputState(unsigned id);
*/
import "C"
import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"unsafe"

	"github.com/lanzafame/bobblehat/sense/screen/color"
	"github.com/tangzero/chip8-emulator/chip8"
)

const (
	RetroButtonB uint8 = iota
	RetroButtonY
	RetroButtonSelect
	RetroButtonStart
	RetroButtonUp
	RetroButtonDown
	RetroButtonLeft
	RetroButtonRight
	RetroButtonA
	RetroButtonX
	RetroButtonL
	RetroButtonR
	FirstRetroButton = RetroButtonB
	LastRetroButton  = RetroButtonR
)

var KeyMapping = map[uint8]uint8{
	RetroButtonLeft:  0x04,
	RetroButtonX:     0x05,
	RetroButtonRight: 0x06,
}

var BuildVersion string
var Emulator *chip8.Emulator
var FrameBuffer *color.RGB565
var KeysState [16]bool

//export Initialize
func Initialize() {
	playSound := func() {}
	stopSound := func() {}
	soundPlayer := func([]byte) (func(), func()) { return playSound, stopSound }

	Emulator = chip8.NewEmulator(KeyPressed, soundPlayer)

	FrameBuffer = color.NewRGB565(Emulator.Display.Rect)
}

//export Deinitialize
func Deinitialize() {
	Emulator = nil
}

//export GetEmulatorInfo
func GetEmulatorInfo(info *C.retro_system_info) {
	info.library_name = C.CString("CHIP-8 Emulator by TangZero")
	info.library_version = C.CString(BuildVersion)
	info.valid_extensions = C.CString("ch8")
	info.need_fullpath = true
	info.block_extract = false
}

//export GetEmulatorAVInfo
func GetEmulatorAVInfo(info *C.retro_system_av_info) {
	info.geometry.base_width = chip8.Width
	info.geometry.base_height = chip8.Height
	info.geometry.max_width = chip8.Width
	info.geometry.max_height = chip8.Height
	info.geometry.aspect_ratio = 0.0
	info.timing.fps = chip8.FPS
	info.timing.sample_rate = 44100
}

//export Reset
func Reset() {
	Emulator.Reset()
}

//export Run
func Run() {
	C.InputPoll()
	UpdateKeysState()

	Emulator.Update()

	// convert from RGBA to RGB565
	draw.Draw(FrameBuffer, Emulator.Display.Rect, Emulator.Display, image.Point{}, draw.Src)

	// draw frame
	C.VideoRefresh(unsafe.Pointer(&FrameBuffer.Pix[0]), chip8.Width, chip8.Height, C.size_t(FrameBuffer.Stride))
}

//export LoadGame
func LoadGame(game *C.retro_game_info) bool {
	data, err := ioutil.ReadFile(C.GoString(game.path))
	if err != nil {
		log.Println(err)
		return false
	}
	Emulator.LoadROM(chip8.ROM{Data: data})
	return true
}

func UpdateKeysState() {
	for button := FirstRetroButton; button < LastRetroButton; button++ {
		key, ok := KeyMapping[button]
		if ok {
			KeysState[key] = C.InputState(C.uint(button)) > 0
		}
	}
}

func KeyPressed(key uint8) bool {
	return KeysState[key]
}

func main() {}
