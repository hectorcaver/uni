/******************************************************************************/
/*                                                                            */
/* project  : PRACTICAS SE-II UNIZAR                                          */
/* filename : main.c                                                         */
/* version  : 3                                                               */
/* date     : 28/09/2020                                                      */
/* author   : Jose Luis Villarroel                                            */
/* description : Medicion de WCET. PR11                                       */
/*                                                                            */
/******************************************************************************/

/******************************************************************************/
/*                        Used modules                                        */
/******************************************************************************/

#include <stdbool.h>

#include <xdc/std.h>

#include <xdc/runtime/System.h>

#include <xdc/runtime/Log.h>
#include <ti/uia/events/UIABenchmark.h>

#include <ti/sysbios/BIOS.h>
#include <ti/sysbios/knl/Clock.h>
#include <ti/sysbios/knl/Task.h>

#include <xdc/runtime/Types.h>

#include "inc/hw_memmap.h"
#include "driverlib/gpio.h"
#include "driverlib/sysctl.h"

#include "computos.h"

/******************************************************************************/
/*                        Global variables                                    */
/******************************************************************************/

Task_Handle task;


/******************************************************************************/
/*                        Tasks                                               */
/******************************************************************************/

Void measure (UArg a0, UArg a1)

{

    for (;;) {

        Log_write1(UIABenchmark_start, (xdc_IArg)"WCET");
        GPIOPinWrite(GPIO_PORTD_BASE, GPIO_PIN_2,GPIO_PIN_2) ;

        CS (20) ;

        GPIOPinWrite(GPIO_PORTD_BASE, GPIO_PIN_2,0) ;
        Log_write1(UIABenchmark_stop, (xdc_IArg)"WCET");

        Task_sleep(2);
    }
}

/******************************************************************************/
/*                        main                                                */
/******************************************************************************/

Void main()
{ 

    System_printf("enter main()\n");

    SysCtlPeripheralEnable(SYSCTL_PERIPH_GPIOD);
    GPIOPinTypeGPIOOutput(GPIO_PORTD_BASE, GPIO_PIN_2);

    task = Task_create(measure, NULL, NULL);

    BIOS_start();     /* enable interrupts and start SYS/BIOS */
}
