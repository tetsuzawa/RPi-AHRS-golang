package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func stringToFloat64(s []string) []float64 {
	f := make([]float64, len(s))
	for n := range s {
		f[n], _ = strconv.ParseFloat(s[n], 64)
	}
	return f
}

func splitDataToFloat64(s, sep string) (float64, float64, float64, float64, float64, float64, float64, float64, float64, error) {
	ssep := strings.Split(s, sep)
	f := stringToFloat64(ssep)
	if len(f) != 9 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, errors.New("received data is incorrect, ")
	}
	return f[0], f[1], f[2],
		f[3], f[4], f[5], f[6], f[7], f[8], nil
}

func updateParams(p *Parameters, serverIp string, sigCh chan os.Signal, stopCh chan struct{}, wg *sync.WaitGroup) {
	//Done
	defer wg.Done()
	// defer func() { wg.Done() }()
	defer fmt.Println("done updateParams")

	// conn, err := net.ListenPacket("udp", "127.0.0.1:62000")
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
		case <-stopCh:
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
			// log.Printf("\nReceived from %v\n", remoteAddr)

			p.wx, p.wy, p.wz, p.ax, p.ay, p.az, p.mx, p.my, p.mz, err = splitDataToFloat64(data, ",")
			if err != nil {
				log.Println("Value Error:", err)
				sigCh <- os.Interrupt
			}
		}
	}
}
