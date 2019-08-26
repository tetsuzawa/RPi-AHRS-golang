package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func sendAttitude(roll, pitch, yaw *float64, stopCh chan struct{}, wg *sync.WaitGroup) {

	//Done
	defer wg.Done()
	// defer func() {wg.Done()}()
	defer log.Println("done sendAttitude")

	conn, err := net.Dial("udp", "127.0.0.1:60002")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var msg = make([]byte, 0, 24)
	// メッセージを送信する
	for {
		// Stop
		select {
		case <- stopCh:
			log.Println("(goroutine sendAttitude) stop request received")
			return
		default:
			// log.Println("(goroutine sendAttitude) runnning")
			time.Sleep(10 * time.Millisecond)
			rollB := float64ToByte(*roll)
			pitchB := float64ToByte(*pitch)
			yawB := float64ToByte(*yaw)
			msg = append(msg, rollB...)
			msg = append(msg, pitchB...)
			msg = append(msg, yawB...)
			conn.Write(msg)
		}
	}
}
