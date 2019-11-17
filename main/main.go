package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/westphae/quaternion"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Parameters struct {
	wx, wy, wz, ax, ay, az, mx, my, mz, dt, beta float64
}

var roll, pitch, yaw float64

func radianToDegree(rad float64) (deg float64) {
	deg = rad * 180 / math.Pi
	return
}

func main() {
	//read ip
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("too many arguments. Usage: %v <your ip>\nfor example %v 172.24.176.10\n", args[0], args[0])
		os.Exit(1)
	} else if len(args) < 2 {
		fmt.Printf("too few arguments. Usage: %v <your ip>\nfor example %v 172.24.176.10\n", args[0], args[0])
		os.Exit(1)
	}
	serverIp := args[1]

	defer fmt.Printf("\nend!!\n")
	defer fmt.Println("done main")

	// declare a channel to receive signal
	sigCh := make(chan os.Signal, 1)
	// receive
	signal.Notify(sigCh, os.Interrupt)
	//pass the channel to the main processing function
	ctx, cancel := context.WithCancel(context.Background())

	// ******** signal handler ********
	wg := sync.WaitGroup{}

	// wait a signal in another gotoutine
	go func() {
		// block until func receive a signal
		sig := <-sigCh
		fmt.Println("Got signal", sig)
		defer cancel()
		return
	}()
	// ********************************

	q := quaternion.Quaternion{W: 1, X: 0, Y: 0, Z: 0}
	roll, pitch, yaw = q.Euler()
	var params = Parameters{wx: 0.01, wy: 0.01, wz: 0.01, ax: 0.01, ay: 0.01, az: 1.00,
		mx: 30, my: 0, mz: 0, dt: 0.02, beta: 1.0,}

	// 通信読込 + 接続相手アドレス情報が受取
	wg.Add(1)
	go updateParams(ctx, &params, serverIp, sigCh, &wg)
	//wg.Add(1)
	//go sendAttitude(&roll, &pitch, &yaw, stopCh, &wg)

	// ******** input handler ********
	offsetCh := make(chan []float64)
	// *******************************
	wg.Add(1)
	go inputHandler(ctx, offsetCh, &wg)


	var st = time.Now()
	var i = 1
	var rollOffset, pitchOffset, yawOffset float64

MainFor:
	for {
		select {
		case <-ctx.Done():
			break MainFor

		case offsets := <-offsetCh:
			rollOffset = offsets[0]
			pitchOffset = offsets[1]
			yawOffset = offsets[2]


		default:
			params.dt = time.Since(st).Seconds()
			// update attitude
			params.updateAttitude(&q)
			st = time.Now()
			roll, pitch, yaw = q.Euler()

			roll = radianToDegree(roll)
			pitch = radianToDegree(pitch)
			yaw = radianToDegree(yaw)

			//log.Printf("%d %3.0f %3.0f %3.0f \r", i, roll, pitch, yaw)
			fmt.Printf("%d roll: %+3.0f, pitch: %+3.0f, yaw: %+3.0f \r", i, roll-rollOffset, pitch-pitchOffset, yaw-yawOffset)

			time.Sleep(10 * time.Millisecond)
			i++
		}
	}

	time.Sleep(1 * time.Second)
	_, err := fmt.Fprintf(os.Stdout, "\n***** Press enter to terminate *****\n")
	check(err)
	wg.Wait()
}

func inputHandler(ctx context.Context, offsetCh chan []float64, wg *sync.WaitGroup) {
	//func inputHandler(stopCh chan struct{}, offsetCh chan []float64, wg *sync.WaitGroup) {
	defer wg.Done()

	var offset []float64

	sc := bufio.NewScanner(os.Stdin)
	for {
		//time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			//case <-stopCh:
			fmt.Println("done inputHandler")
			return

		default:
			sc.Scan()
			s := sc.Text()
			switch s {
			case "n":
				fmt.Printf("\nnew offset have been set with values roll: %+3.0f, pitch: %+3.0f, yaw: %+3.0f\n", roll, pitch, yaw)
				offset = []float64{roll, pitch, yaw}
				offsetCh <- offset
			case "s":
				fmt.Printf("\noffset have been set with values roll: %+3.0f, pitch: %+3.0f, yaw: %+3.0f\n", roll, pitch, yaw)
				offsetCh <- offset
			case "c":
				fmt.Printf("\noffset have been canceled\n")
				offsetCh <- []float64{0, 0, 0}
			default:
			}
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
