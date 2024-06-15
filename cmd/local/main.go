package main

import (
	"sync"

	"github.com/chack93/dmx_light_ctrl/internal/service/config"
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
	if err := server.New().Init(wg); err != nil {
		log.Fatalf("server init failed, err: %v", err)
	}

	wg.Wait()
}
