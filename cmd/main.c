// #include <stdio.h>
#include "lib/demo_svc/demo_svc.h"
#include "lib/dmx_svc/dmx_svc.h"
#include "pico/cyw43_arch.h"
// #include "pico/multicore.h"
#include "pico/stdlib.h"

/*
void core1_entry()
{
    DmxSvc_startLoop();
    while (1) {
        tight_loop_contents();
    }
}
*/

#include "hardware/spi.h"
void writeAddress()
{
    // 1/(4*10^-6)
    spi_init(spi_default, 9600);
    gpio_set_function(PICO_DEFAULT_SPI_RX_PIN, GPIO_FUNC_SPI);
    gpio_set_function(PICO_DEFAULT_SPI_SCK_PIN, GPIO_FUNC_SPI);
    gpio_set_function(PICO_DEFAULT_SPI_TX_PIN, GPIO_FUNC_SPI);
    gpio_set_function(PICO_DEFAULT_SPI_CSN_PIN, GPIO_FUNC_SPI);
    uint16_t addr16 = 0x04;
    spi_write16_blocking(spi_default, &addr16, 1);
}

int main()
{
    stdio_init_all();
    if (cyw43_arch_init()) {
        // printf("Wi-Fi init failed");
        return -1;
    }
    writeAddress();
    DemoSvc_init();
    DmxSvc_init();

    /*
            gpio_init(26);
            gpio_set_dir(26, GPIO_OUT);
            gpio_put(26, true);
    */

    DmxSvc_startLoop();
    // multicore_launch_core1(core1_entry);
    while (true) {
        tight_loop_contents();
    }

    /*
while (true) {
    sleep_ms(1000);
    cyw43_arch_gpio_put(CYW43_WL_GPIO_LED_PIN, 1);
    sleep_us(500000);
    cyw43_arch_gpio_put(CYW43_WL_GPIO_LED_PIN, 0);
    sleep_us(1500000);
}
    */
}
