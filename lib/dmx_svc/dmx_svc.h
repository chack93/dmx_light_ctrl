// dmx_svc.h
#ifndef DMX_SVC_H
#define DMX_SVC_H

#include "pico/stdlib.h"
#include "pico/util/queue.h"

typedef struct {
  //uint8_t ledNr;
  uint8_t red;
  uint8_t green;
  uint8_t blue;
  //uint8_t white;
} DmxSvc_lightDevice;

// (512 dmx-channels - 1 start-channel) / 4 led-colors = 127,75 => 127
#define DMX_SVC_DEVICE_LENGTH 127
extern DmxSvc_lightDevice DmxSvc_list[DMX_SVC_DEVICE_LENGTH];
extern DmxSvc_lightDevice DmxSvc_listBuffer[DMX_SVC_DEVICE_LENGTH];

// Initialize dmx light
void DmxSvc_init();
void DmxSvc_startLoop();
int DmxSvc_setLedColor(int led, uint8_t red, uint8_t green, uint8_t blue);

#endif // DMX_SVC_H
