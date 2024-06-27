package dmx_svc

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type DmxSvcColorMsg struct {
	Red   uint8
	Green uint8
	Blue  uint8
}
type DmxSvcSetChannelMsg func(channel int, msg uint8)
type DmxSvcWriteLogicLevel func(high bool, debugMsg string)
type DmxSvcWriteConfig struct {
	Rate int64
}
type MockDataItem struct {
	Time     int64  `json:"time"`
	Data     int    `json:"data"`
	DebugMsg string `json:"dbgMsg"`
}
type DmxSvc struct {
	mockMode         bool
	mockStart        time.Time
	mockData         []MockDataItem
	channel          []uint8
	protocol         DmxSvcImpl
	writeLoopRunning bool
	a                rpio.Pin
	b                rpio.Pin
}

type DmxSvcImpl interface {
	GetConfig() DmxSvcWriteConfig
	InitChannel(setChFn DmxSvcSetChannelMsg)
	SetDeviceColor(deviceNr int, msg DmxSvcColorMsg, setChFn DmxSvcSetChannelMsg)
	WriteDmxMsg(chMsgList []uint8, write DmxSvcWriteLogicLevel)
}
