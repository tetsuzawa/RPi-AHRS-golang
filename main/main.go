package main

import (
	"context"
	"github.com/westphae/quaternion"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Parameters struct {
	wx, wy, wz, ax, ay, az, mx, my, mz, dt, beta float64
}

func radianToDegree(rad float64) (deg float64) {
	deg = rad * 180 / math.Pi
	return
}

func main() {
	defer log.Println("done main")

	// declare signals to treat
	trapSignals := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	// declare a channel to receive signal
	sigCh := make(chan os.Signal, 1)
	// receive
	// signal.Notify(sigCh, trapSignals...)
	signal.Notify(sigCh, os.Interrupt)
	log.Println(trapSignals)
	//pass the channel to the main processing function
	ctx, cancel := context.WithCancel(context.Background())

	// ##################################
	// signal handler
	stopCh := make(chan struct{})
	wg := sync.WaitGroup{}
	// ##################################

	// wait a signal in another gotoutine
	go func() {
		// block until func receive a signal
		sig := <-sigCh
		close(stopCh)
		log.Println("Got signal", sig)
		defer cancel()
		return
	}()

	q := quaternion.Quaternion{W: 1, X: 0, Y: 0, Z: 0}
	roll, pitch, yaw := q.Euler()
	var params = Parameters{
		wx:   0.01,
		wy:   0.01,
		wz:   0.01,
		ax:   0.01,
		ay:   0.01,
		az:   1.00,
		mx:   30,
		my:   0,
		mz:   0,
		dt:   0.02,
		beta: 1.0,
	}
	// 通信読込 + 接続相手アドレス情報が受取
	wg.Add(1)
	go updateParams(&params, sigCh, stopCh, &wg)
	wg.Add(1)
	go sendAttitude(&roll, &pitch, &yaw, stopCh, &wg)

	st := time.Now()
	i := 1

MainFor:
	for {
		select {
		case <-ctx.Done():
			break MainFor
		default:
			params.dt = time.Since(st).Seconds()
			// update attitude
			params.updateAttitude(&q)
			st = time.Now()
			roll, pitch, yaw = q.Euler()

			roll = radianToDegree(roll)
			pitch = radianToDegree(pitch)
			yaw = radianToDegree(yaw)

			log.Println(i, roll, pitch, yaw)

			time.Sleep(10 * time.Millisecond)
			i++
		}
	}
	wg.Wait()
}
