package main

import (
	"github.com/google/periph/experimental/devices/mpu9250"
	"log"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
	"time"
)

func main() {
	host.Init()
	log.Println("start mpu")
	setMPU9250()
	log.Println("start led")
	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		rpi.P1_33.Out(l)
		<-t.C
	}
}

func setMPU9250() {
	var ahrs = mpu9250.MPU9250{}
	ahrs.Calibrate()
	// var acc *mpu9250.AccelerometerData
	acc, err := ahrs.GetAcceleration()
	if err != nil {
		log.Fatalln(err)
	}
	println(acc.X, acc.Y, acc.Z)

}
