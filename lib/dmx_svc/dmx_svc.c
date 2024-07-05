#include "dmx_svc.h"
#include "impl_tm512ac.h"
#include "pico/cyw43_arch.h"
#include "pico/stdlib.h"
//#include "pico/util/queue.h"

DmxSvc_lightDevice DmxSvc_list[DMX_SVC_DEVICE_LENGTH];
DmxSvc_lightDevice DmxSvc_listBuffer[DMX_SVC_DEVICE_LENGTH];
//queue_t dmxSvc_listSendQueue;

void copyBuffer()
{
    /*
    int count = queue_get_level(&dmxSvc_listSendQueue);
    if (count) {
        for (; count > 0; count--) {
            DmxSvc_lightDevice el;
            queue_remove_blocking(&dmxSvc_listSendQueue, &el);
            DmxSvc_list[el.ledNr] = el;
        }
    }
    */
    for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
        DmxSvc_list[i] = DmxSvc_listBuffer[i];
    }
}
void DmxSvc_init()
{
    //queue_init(&dmxSvc_listSendQueue, sizeof(DmxSvc_lightDevice), DMX_SVC_DEVICE_LENGTH);
    for (int i = 0; i < DMX_SVC_DEVICE_LENGTH; i++) {
        //DmxSvc_list[i].ledNr = i;
        DmxSvc_list[i].red = 0;
        DmxSvc_list[i].green = 0;
        DmxSvc_list[i].blue = 0;
        DmxSvc_listBuffer[i].red = 0;
        DmxSvc_listBuffer[i].green = 0;
        DmxSvc_listBuffer[i].blue = 0;
    }
}
void DmxSvc_startLoop()
{
    dmxSvcImpl_tm512acInit();
    int i = 0;
    while (true) {
        dmxSvcImpl_tm512acWrite();
        copyBuffer();
        sleep_ms(10);
        // tight_loop_contents();
    }
}
int DmxSvc_setLedColor(int led, uint8_t red, uint8_t green, uint8_t blue)
{
    if (led < 0 || led >= DMX_SVC_DEVICE_LENGTH) {
        return -1;
    }

    DmxSvc_lightDevice el = {
        //.ledNr = led,
        .red = red,
        .green = green,
        .blue = blue,
        //.white = white
    };
    //queue_try_add(&dmxSvc_listSendQueue, &el);

    DmxSvc_listBuffer[led].red = red;
    DmxSvc_listBuffer[led].green = green;
    DmxSvc_listBuffer[led].blue = blue;
    return 0;
}
