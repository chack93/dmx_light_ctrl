package dmx_svc

import (
	"fmt"
	"time"
)

const CH_START = 1
const CH_RED = 0
const CH_GREEN = 1
const CH_BLUE = 2
const CH_INTERVAL = 4
const CH_RATE = 250_000

type TM512AC struct{}

func (s *TM512AC) GetConfig() DmxSvcWriteConfig {
	return DmxSvcWriteConfig{
		Rate: CH_RATE,
	}
}
func (s *TM512AC) InitChannel(setChFn DmxSvcSetChannelMsg) {
	for i := 0; i < CH_START; i++ {
		setChFn(i, 0x00)
	}
}

func (s *TM512AC) SetDeviceColor(deviceNr int, msg DmxSvcColorMsg, setChFn DmxSvcSetChannelMsg) {
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
	// tm512ac max possible device/channel: 127/509 (1+4*127 = 509 channel)
	chStart := CH_START + deviceNr*CH_INTERVAL
	setChFn(chStart+CH_RED, msg.Red)
	setChFn(chStart+CH_GREEN, msg.Green)
	setChFn(chStart+CH_BLUE, msg.Blue)
}
func (s *TM512AC) WriteDmxMsg(chMsgList []uint8, write DmxSvcWriteLogicLevel) {
	// break signal between: 88µs - 1sec
	write(false, "break")
	time.Sleep(time.Microsecond * 88 * 2)

	// mark after break between: 8µs - 1sec
	write(true, "mab")
	time.Sleep(time.Microsecond * 8)

	// write all channels
	// TODO: maybe write MSB first
	const waitTime = time.Microsecond * time.Duration(1/CH_RATE)
	for chIdx, chMsg := range chMsgList {
		// a channel bit between: 1µs - 5µs

		// channel start bit
		write(false, fmt.Sprintf("ch:%d-S", chIdx))
		time.Sleep(waitTime)

		// write each bit
		for bit := uint8(8); bit > 0; bit-- {
			write(chMsg&(bit-1) == 0x01, fmt.Sprintf("ch:%d-%d", chIdx, bit-1))
			time.Sleep(waitTime)
		}

		// 2 channel stop bits
		write(true, fmt.Sprintf("ch:%d-T", chIdx))
		time.Sleep(waitTime * 2)

		// wait time between fields: 0µs - 1sec
		// time.Sleep(time.Microsecond * 4)
	}
}
