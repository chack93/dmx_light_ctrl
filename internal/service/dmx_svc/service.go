package dmx_svc

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chack93/dmx_light_ctrl/internal/service/logger"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
)

var log = logger.Get()
var svc *DmxSvc

const MOCK_DURATION_MICROSECONDS = 100_000

//go:embed mock-diagram.html
var mockDiagramHtml []byte

func Get() *DmxSvc {
	if svc == nil {
		svc = &DmxSvc{}
	}
	return svc
}

func (svc *DmxSvc) commonInit() error {

	svc.channel = make([]uint8, viper.GetInt("dmx.channel.total"))

	switch viper.GetString("dmx.protocol") {
	case "tm512ac":
		svc.protocol = &TM512AC{}
	default:
		log.Errorf("dmx_svc/commonInit - failed, unknown protocol: %s", viper.GetString("dmx.protocol"))
		return errors.New("unknown dmx.protocol")
	}

	svc.writeLoopRunning = true
	go svc.writeLoop()

	return nil
}

func (svc *DmxSvc) InitMock(wg *sync.WaitGroup) error {
	defer wg.Done()

	svc.mockMode = true
	svc.mockStart = time.Now()
	time.AfterFunc(time.Microsecond*MOCK_DURATION_MICROSECONDS, func() {
		//svc.writeLoopRunning = false
		writeMockDiagram()
	})
	if err := svc.commonInit(); err != nil {
		return err
	}

	return nil
}

func (svc *DmxSvc) Init(wg *sync.WaitGroup) error {
	defer wg.Done()

	svc.mockMode = false
	if err := rpio.Open(); err != nil {
		log.Errorf("dmx_svc/Init - rpio open memory failed: %v", err)
		return err
	}
	svc.a = rpio.Pin(viper.GetInt("dmx.pin.a"))
	svc.b = rpio.Pin(viper.GetInt("dmx.pin.b"))
	svc.a.Output()
	svc.b.Output()
	if err := svc.commonInit(); err != nil {
		return err
	}

	return nil
}

func (svc *DmxSvc) writeLoop() {
	time.Sleep(time.Millisecond * 100)

	/*
		log.Debugf("dmx_svc/writeLoop - channel-length: %d, msg: %x %x,%x,%x %x,%x,%x",
			len(svc.channel),
			svc.channel[0],
			svc.channel[1],
			svc.channel[2],
			svc.channel[3],
			svc.channel[4],
			svc.channel[5],
			svc.channel[6],
		)
	*/
	svc.protocol.WriteDmxMsg(svc.channel, svc.writeLogicLevel)

	if svc.writeLoopRunning {
		svc.writeLoop()
	}
}

func (svc *DmxSvc) writeLogicLevel(high bool, debugMsg string) {
	if svc.mockMode {
		time := time.UnixMicro(0).Add(time.Now().Sub(svc.mockStart)).UnixMicro()
		level := 0
		if high {
			level = 1
		}
		svc.mockData = append(svc.mockData, MockDataItem{Time: time, Data: level, DebugMsg: debugMsg})
	} else {
		if high {
			svc.b.Low()
			svc.a.High()
		} else {
			svc.b.High()
			svc.a.Low()
		}
	}
}

func (svc *DmxSvc) SetDeviceColor(deviceNr int, msg DmxSvcColorMsg) {
	//log.Infof("dmx_svc/SetDeviceColor - deviceNr: %d, msg r,g,b: %d,%d,%d", deviceNr, msg.Red, msg.Green, msg.Blue)
	svc.protocol.SetDeviceColor(deviceNr, msg, svc.setChannelMsg)
}

func (svc *DmxSvc) setChannelMsg(channel int, msg uint8) {
	// just ignore wrong addressing
	if channel >= viper.GetInt("dmx.channel.total") {
		return
	}
	svc.channel[channel] = msg
	//log.Infof("dmx_svc/setChannelMsg - channel: %d, msg: 0x%x", channel, msg)
}

func writeMockDiagram() {
	// csv
	csvFile, err := os.Create("mock-data.csv")
	if err != nil {
		log.Errorf("dmx_svc/writeMockDiagram - failed to create csv, err: %v", err)
	}
	csvFile.WriteString("Time[s], level\n")
	for idx, el := range svc.mockData {
		var lastTime int64 = 0
		if idx != 0 {
			lastTime = svc.mockData[idx-1].Time
		}
		duration := el.Time - lastTime
		csvFile.WriteString(fmt.Sprintf("%d,\t%d,\t%d,\t%s\n", el.Time, duration, el.Data, el.DebugMsg))
	}
	csvFile.Close()

	// custom diagram renderer
	jsonData, err := json.Marshal(svc.mockData)
	if err != nil {
		log.Errorf("dmx_svc/writeMockDiagram - failed to marshal mock data json, err: %v", err)
	}
	mockDiagramHtmlString := string(mockDiagramHtml)
	mockDiagramPage := strings.ReplaceAll(mockDiagramHtmlString, "{JSON_DATA}", string(jsonData))
	diagramPageFile, err := os.Create("mock-diagram.html")
	if err != nil {
		log.Errorf("dmx_svc/writeMockDiagram - failed to create mock diagram page, err: %v", err)
	}
	diagramPageFile.WriteString(mockDiagramPage)
	csvFile.Close()
}
