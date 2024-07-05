// impl_tm512ac.h
#ifndef DMX_SVC_TM512AC_H
#define DMX_SVC_TM512AC_H

#include "dmx_svc.h"
#include "pico/stdlib.h"

#define DMX_SVC_IMPL_A_GPIO_PIN 16
#define DMX_SVC_IMPL_B_GPIO_PIN 17

void writeLvl(bool high);
inline void wait(uint64_t delay_us);
void writeChannel(uint8_t chMsg, int waitTime);

int dmxSvcImpl_tm512acInit()
{
	gpio_init(DMX_SVC_IMPL_A_GPIO_PIN);
	gpio_init(DMX_SVC_IMPL_B_GPIO_PIN);
	gpio_set_dir(DMX_SVC_IMPL_A_GPIO_PIN, GPIO_OUT);
	gpio_set_dir(DMX_SVC_IMPL_B_GPIO_PIN, GPIO_OUT);
	writeLvl(true);
	return 0;
}

void writeLvl(bool high)
{
	if (high) {
		//gpio_put_masked(uint32_t mask, uint32_t value)
		gpio_put(DMX_SVC_IMPL_A_GPIO_PIN, 1);
		gpio_put(DMX_SVC_IMPL_B_GPIO_PIN, 0);
	} else {
		gpio_put(DMX_SVC_IMPL_A_GPIO_PIN, 0);
		gpio_put(DMX_SVC_IMPL_B_GPIO_PIN, 1);
	}
}

inline void wait(uint64_t delay_us)
{
	busy_wait_us(delay_us);
}

void writeChannel(uint8_t chMsg, int waitTime)
{

	// channel start bit
	writeLvl(false);
	wait(waitTime);

	// write each bit
	/*
	//MSB first
	for (int bit = 8; bit > 0; bit--) {
		writeLvl((chMsg & (1 << (bit - 1))) > 0);
		wait(waitTime);
	}
	*/
	// LSB first
	for (int bit = 0; bit < 8; bit++) {
		writeLvl((chMsg & (1 << bit)) > 0);
		wait(waitTime);
	}

	// 2 channel wide stop bits
	writeLvl(true);
	wait(waitTime * 2);

	// wait time between fields: 0µs - 1sec
	// wait(333);
}

#include "hardware/spi.h"
void dmxSvcImpl_tm512acWrite()
{
	// tm512ac dmx msg structure:
	// ch 0: 0x00
	// ch 1: device-1 red
	// ch 2: device-1 green
	// ch 3: device-1 blue
	// ch 4: device-1 white
	// ch 5: device-2 red
	// ch 6: device-2 green
	// ch 7: device-2 blue
	// ch 8: device-2 white
	// etc.
	// tm512ac max possible device/channel: 127/509 => (1+4*127 = 509 channel)

	// break signal between: 88µs - 1sec
	writeLvl(false);
	wait(88);
    uint16_t addr16 = 0x00;
    spi_write16_blocking(spi_default, &addr16, 1);

	// mark after break between: 8µs - 1sec
	writeLvl(true);
	wait(8);

	// a channel bit between: 1µs - 5µs
	// 1/(4*10^-6)
	const int waitTime = 4;

	// first channel 0x00
	writeLvl(false);
	wait(waitTime * 9);
	writeLvl(true);
	wait(waitTime * 2);

	// write all channels
	for (int devIdx; devIdx < DMX_SVC_DEVICE_LENGTH; devIdx++) {
		writeChannel(DmxSvc_list[devIdx].red, waitTime);
		writeChannel(DmxSvc_list[devIdx].green, waitTime);
		writeChannel(DmxSvc_list[devIdx].blue, waitTime);
		//writeChannel(dmxSvc_list[devIdx].white, waitTime);
	}
}

#endif // DMX_SVC_TM512AC_H
