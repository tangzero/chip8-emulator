TARGET=chip8_libretro
LIBRETRO_CORE=$(TARGET).dylib
LIBRETRO_HEADER=$(TARGET).h
BUILD_VERSION=$(shell git rev-parse --short HEAD)

# run: $(LIBRETRO_CORE)
# 	retroarch -v -L $(LIBRETRO_CORE) roms/test_opcode.ch8

$(LIBRETRO_CORE): clean
	go build -buildmode=c-shared -ldflags "-X 'main.BuildVersion=$(BUILD_VERSION)'" -o $@ ./libretro

clean: 
	rm -f $(LIBRETRO_CORE) $(LIBRETRO_HEADER)

test:
	go test -v -race ./chip8

