package main

/*
#cgo CFLAGS: -I./libretro-common/include
#include "libretro.h"

typedef struct retro_system_info retro_system_info;
typedef struct retro_system_av_info retro_system_av_info;
*/
import "C"
import (
	"fmt"
)

const RetroAPIVersion = 1

//export retro_init
func retro_init() {
	fmt.Println("loading chip8 core...")
}

//export retro_deinit
func retro_deinit() {
	fmt.Println("unloading chip8 core...")
}

//export retro_api_version
func retro_api_version() uint32 {
	return RetroAPIVersion
}

func main() {}
