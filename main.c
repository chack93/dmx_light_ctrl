//#include <stdio.h>
#include "hardware/pio.h"
#include "hardware/timer.h"
#include "pico/cyw43_arch.h"
#include "pico/stdlib.h"
#include "pico/cyw43_arch.h"

int main() {
    stdio_init_all();
    if (cyw43_arch_init()) {
        //printf("Wi-Fi init failed");
        return -1;
    }
    while (true) {
        cyw43_arch_gpio_put(CYW43_WL_GPIO_LED_PIN, 1);
        sleep_us(500000);
        cyw43_arch_gpio_put(CYW43_WL_GPIO_LED_PIN, 0);
        sleep_us(1500000);
    }
}


// TODO:
// 1. make channel list 512
// 2. led demo function
// 3. dmx channel worker
// 4. dmx protocol writer

/*
int64_t alarm_callback(alarm_id_t id, void *user_data) {
    // Put your timeout handler code in here
    return 0;
}



int main()
{
    stdio_init_all();


    // Timer example code - This example fires off the callback after 2000ms
    add_alarm_in_ms(2000, alarm_callback, NULL, false);


    puts("Hello, world!");

    return 0;
}
*/
