TARGET=chip8_libretro
LIBRETRO_CORE=$(TARGET).so
LIBRETRO_HEADER=$(TARGET).h

# run: $(LIBRETRO_CORE)
# 	retroarch -v -L $(LIBRETRO_CORE) roms/test_opcode.ch8

$(LIBRETRO_CORE):
	go build -buildmode=c-shared -o $@ ./libretro

test:
	go test -v -race ./chip8

clean: 
	rm -f $(LIBRETRO_CORE) $(LIBRETRO_HEADER)

