package main

/*
#cgo CFLAGS: -I./libretro-common/include
#include "libretro.h"

typedef struct retro_system_info retro_system_info;
typedef struct retro_system_av_info retro_system_av_info;
*/
import "C"
import (
	"github.com/tangzero/chip8-emulator/chip8"
)

var BuildVersion string
var emulator *chip8.Emulator

//export Initialize
func Initialize() {
	emulator = chip8.NewEmulator()
}

//export Deinitialize
func Deinitialize() {
	emulator = nil
}

//export GetEmulatorInfo
func GetEmulatorInfo(info *C.retro_system_info) {
	info.library_name = C.CString("CHIP-8 Emulator by TangZero")
	info.library_version = C.CString(BuildVersion)
	info.valid_extensions = C.CString("ch8")
	info.need_fullpath = false
	info.block_extract = false
}

//export GetEmulatorAVInfo
func GetEmulatorAVInfo(info *C.retro_system_av_info) {
	info.geometry.base_width = 64
	info.geometry.base_height = 32
	info.geometry.max_width = 64
	info.geometry.max_height = 32
	info.geometry.aspect_ratio = 0.0
	info.timing.fps = 60
	info.timing.sample_rate = 44100
}

func main() {}
