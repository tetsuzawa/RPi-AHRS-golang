package main

import (
	"fmt"
	"log"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
	"time"
)

func main() {
	fmt.Println("On Raspberry Pi")
	log.Println("MPU9250")
	host.Init()
	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		rpi.P1_33.Out(l)
		<-t.C
	}
}
