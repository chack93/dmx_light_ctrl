#include <stdint.h>
#include <stdio.h>

#define DMX_SVC_DEVICE_LENGTH 36
#define LEVEL 1.0
typedef struct {
  uint8_t red;
  uint8_t green;
  uint8_t blue;
  uint8_t white;
} DmxSvc_lightDevice;
extern DmxSvc_lightDevice dmxSvc_list[];

int testFnStatic();
int testFnRainbow();
int main()
{

  testFnStatic();

  return 0;
  for (int i = 0; i < 255; i++) {
    testFnRainbow();
  }
  return 0;
}

int DmxSvc_setLedColor(int led, uint8_t red, uint8_t green, uint8_t blue, uint8_t white)
{
  printf("led: %d, %02x%02x%02x\n", led, red, green, blue);
  return 0;
}

static int demoSvc_colorIdx = 0;
int testFnStatic()
{
  for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
    uint8_t red = (i  + demoSvc_colorIdx) % 3 == 0 ? 255 : 0;
    uint8_t green = (i  + demoSvc_colorIdx) % 3 == 1 ? 255 : 0;
    uint8_t blue = (i  + demoSvc_colorIdx) % 3 == 2 ? 255 : 0;

    red *= LEVEL;
    green *= LEVEL;
    blue *= LEVEL;
    uint8_t white = 0;
    DmxSvc_setLedColor(i, red, green, blue, white);
  }
  demoSvc_colorIdx = demoSvc_colorIdx + 1;
  return 0;
}

int testFnRainbow()
{
  for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
    const uint8_t third = 255 / 3;

    uint8_t red = demoSvc_colorIdx + (third * (i % 3));
    uint8_t green = demoSvc_colorIdx + (third * (i % 3));
    uint8_t blue = demoSvc_colorIdx + (third * (i % 3));

    red *= LEVEL;
    green *= LEVEL;
    blue *= LEVEL;
    uint8_t white = 0;
    DmxSvc_setLedColor(i, red, green, blue, white);
  }
  demoSvc_colorIdx = demoSvc_colorIdx + 1;
  return 0;
}
