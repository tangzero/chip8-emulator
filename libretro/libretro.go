package main

/*
#include "libretro-common/include/libretro.h"

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
	emulator = chip8.NewEmulator(nil)
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
	info.geometry.base_width = chip8.Width
	info.geometry.base_height = chip8.Height
	info.geometry.max_width = chip8.Width
	info.geometry.max_height = chip8.Height
	info.geometry.aspect_ratio = 0.0
	info.timing.fps = chip8.FPS
	info.timing.sample_rate = 44100
}

func main() {}
