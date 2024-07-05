APP_NAME = dmx_light_ctrl
VERSION = 1.0.0

.PHONY: help
help:
	@echo "make options\n\
		- all             clean, cmake, build, flash\n\
		- clean           clean build directory bin/\n\
		- cmake           run cmake & create build dir\n\
		- cmake-dbg       run cmake with debug flag & create build dir\n\
		- build           build binary\n\
		- flash           copy built binary to /Volumes/RPI-RP2\n\
		- deploy          build & flash\n\
		- deploy-dbg      build debug bin & flash\n\
		- help            display this message"

.PHONY: all
all: clean cmake build flash

.PHONY: clean
clean:
	rm -rf build

.PHONY: cmake
cmake: clean
	mkdir build
	(cd build && cmake ..)

.PHONY: cmake-dbg
cmake-dbg: clean
	mkdir build
	(cd build && cmake -DCMAKE_BUILD_TYPE=Debug ..)

.PHONY: build
build:
	(cd build && make)

.PHONY: flash
flash: |
	cp build/dmx_light_ctrl.uf2 /Volumes/RPI-RP2/. || \
		sudo openocd -f interface/cmsis-dap.cfg -f target/rp2040.cfg -c "adapter speed 5000" -c "program build/dmx_light_ctrl.elf verify reset exit"

.PHONY: deploy
deploy: build flash

.PHONY: deploy-dbg
deploy-dbg: cmake-dbg build flash
	sudo openocd -f interface/cmsis-dap.cfg -f target/rp2040.cfg -c "adapter speed 5000"
	# gdb build/dmx_light_ctrl.elf
	# sudo openocd -f interface/cmsis-dap.cfg -f target/rp2040.cfg -c "adapter speed 5000"
	# target remote localhost:3333
	# monitor reset init
	# continue
