package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func updateParams(ctx context.Context, p *Parameters, serverIp string, sigCh chan os.Signal, wg *sync.WaitGroup) {
	//Done
	defer wg.Done()
	defer fmt.Println("done updateParams")

	serverAddr := serverIp + ":50020"
	fmt.Printf("udp server is runnning at %s\n", serverAddr)
	conn, err := net.ListenPacket("udp", serverAddr)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	var data string

	for {
		// Stop
		select {
		case <-ctx.Done():
			fmt.Println("(goroutine updateParams) stop request received")
			return
		default:
			time.Sleep(5 * time.Millisecond)
			if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
				log.Println(err)
			}
			length, _, err := conn.ReadFrom(buffer)
			if err != nil {
				log.Println("Connect ERROR : ", err)
				continue
			}
			data = string(buffer[:length])

			p.wx, p.wy, p.wz, p.ax, p.ay, p.az, p.mx, p.my, p.mz, err = splitDataToFloat64(data, ",")
			if err != nil {
				log.Println("Value Error:", err)
				sigCh <- os.Interrupt
			}
		}
	}
}
