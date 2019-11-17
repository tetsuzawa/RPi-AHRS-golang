package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
)

func offsetHandler(ctx context.Context, offsetCh chan []float64, wg *sync.WaitGroup) {
	//func offsetHandler(stopCh chan struct{}, offsetCh chan []float64, wg *sync.WaitGroup) {
	defer wg.Done()

	var offset []float64

	sc := bufio.NewScanner(os.Stdin)
	for {
		//time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			//case <-stopCh:
			fmt.Println("done offsetHandler")
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
