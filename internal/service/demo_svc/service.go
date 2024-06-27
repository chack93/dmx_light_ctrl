package demo_svc

import (
	"sync"
	"time"

	"github.com/chack93/dmx_light_ctrl/internal/service/dmx_svc"
)

type DemoLight int

const (
	DemoLightNone DemoLight = iota
	DemoLightRed
	DemoLightGreen
	DemoLightBlue
	DemoLightRainbow
	DemoLightWave
)

var demoLoop = []DemoLight{
	DemoLightRed,
	DemoLightGreen,
	DemoLightBlue,
	DemoLightRainbow,
	DemoLightWave,
}

type DemoSvc struct {
	Current DemoLight
}

var svc *DemoSvc

func Get() *DemoSvc {
	if svc == nil {
		svc = &DemoSvc{}
	}
	return svc
}

func (svc *DemoSvc) Init(wg *sync.WaitGroup) error {
	defer wg.Done()
	svc.Current = DemoLightRed

	go svc.RunDemoLoop()
	//dmx_svc.Get().Channel

	return nil
}

func (svc *DemoSvc) RunDemoLoop() {

	dim := 0.5
	a := uint8(254.0 * dim)
	b := uint8(127.0 * dim)
	c := uint8(63.0 * dim)

	for i := 0; i < 3; i++ {
		dmx_svc.Get().SetDeviceColor((i+0)%3, dmx_svc.DmxSvcColorMsg{Red: uint8(a), Green: uint8(b), Blue: uint8(c)})
		dmx_svc.Get().SetDeviceColor((i+1)%3, dmx_svc.DmxSvcColorMsg{Red: uint8(b), Green: uint8(c), Blue: uint8(a)})
		dmx_svc.Get().SetDeviceColor((i+2)%3, dmx_svc.DmxSvcColorMsg{Red: uint8(c), Green: uint8(a), Blue: uint8(b)})
		time.Sleep(time.Millisecond * 5000)
	}
	if svc.Current != DemoLightNone {
		svc.RunDemoLoop()
	}
}
func (svc *DemoSvc) StopDemoLoop() {
	svc.Current = DemoLightNone
}
