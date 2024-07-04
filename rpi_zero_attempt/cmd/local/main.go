package main

import (
	"sync"

	"github.com/chack93/dmx_light_ctrl/internal/service/config"
	"github.com/chack93/dmx_light_ctrl/internal/service/demo_svc"
	"github.com/chack93/dmx_light_ctrl/internal/service/dmx_svc"
	"github.com/chack93/dmx_light_ctrl/internal/service/logger"
	"github.com/chack93/dmx_light_ctrl/internal/service/server"
)

var log = logger.Get()

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("config init failed, err: %v", err)
	}
	if err := logger.Init(); err != nil {
		log.Fatalf("log init failed, err: %v", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	if err := dmx_svc.Get().InitMock(wg); err != nil {
		log.Fatalf("dmx svc mock init failed, err: %v", err)
	}
	wg.Add(1)
	if err := demo_svc.Get().Init(wg); err != nil {
		log.Fatalf("demo svc init failed, err: %v", err)
	}
	wg.Add(1)
	if err := server.New().Init(wg); err != nil {
		log.Fatalf("server init failed, err: %v", err)
	}

	wg.Wait()
}
