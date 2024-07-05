#include "demo_svc.h"
#include "../dmx_svc/dmx_svc.h"
#include "pico/stdlib.h"

enum DemoSvc_ProgramEnum {
    DemoSvc_Program_1,
    DemoSvc_Program_2,
    DemoSvc_Program_3
};
const int demoSvc_program = DemoSvc_Program_2;

static uint demoSvc_colorIdx = 0;
static bool demoSvc_timerRunning = true;
static struct repeating_timer demoSvc_changeLightTimer;

bool change_light_t_callback(struct repeating_timer* rt);
void change_color();

void DemoSvc_init()
{
    add_repeating_timer_ms(1000, change_light_t_callback, NULL, &demoSvc_changeLightTimer);
    change_color();
}

bool change_light_t_callback(struct repeating_timer* rt)
{
    change_color();
    return true;
}

void change_color()
{
    const float level = 0.2;

    switch (demoSvc_program) {
    case DemoSvc_Program_1: {
        if (demoSvc_colorIdx > 1) {
            return;
        }
        uint8_t demoList[] = { 255, 0, 0 };
        uint8_t demoLen = 3;
        for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
            uint8_t red = demoList[(i + demoSvc_colorIdx + 0) % demoLen] * level;
            uint8_t green = demoList[(i + demoSvc_colorIdx + 1) % demoLen] * level;
            uint8_t blue = demoList[(i + demoSvc_colorIdx + 2) % demoLen] * level;
            red *= level;
            green *= level;
            blue *= level;
            DmxSvc_setLedColor(i, red, green, blue);
        }
        break;
    }
    case DemoSvc_Program_2: {
        uint8_t demoList[] = { 255, 0, 0 };
        uint8_t demoLen = 3;
        for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
            uint8_t red = (i + demoSvc_colorIdx) % 3 == 0 ? 255 : 0;
            uint8_t green = (i + demoSvc_colorIdx) % 3 == 1 ? 255 : 0;
            uint8_t blue = (i + demoSvc_colorIdx) % 3 == 2 ? 255 : 0;
            red *= level;
            green *= level;
            blue *= level;
            DmxSvc_setLedColor(i, red, green, blue);
        }
        break;
    }
    case DemoSvc_Program_3: {
        const uint8_t third = 255 / 3;
        for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
            uint8_t red = demoSvc_colorIdx + (1 * third * (i % 3));
            uint8_t green = demoSvc_colorIdx + (2 * third * (i % 3));
            uint8_t blue = demoSvc_colorIdx + (3 * third * (i % 3));
            red *= level;
            green *= level;
            blue *= level;
            DmxSvc_setLedColor(i, red, green, blue);
        }
        demoSvc_colorIdx = demoSvc_colorIdx + 15;
        break;
    }
    }

    demoSvc_colorIdx = demoSvc_colorIdx + 1;
}
