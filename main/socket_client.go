package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func sendAttitude(roll, pitch, yaw *float64, stopCh chan struct{}, wg *sync.WaitGroup) {

	//Done
	defer wg.Done()
	defer log.Println("done sendAttitude")

	conn, err := net.Dial("udp", "127.0.0.1:60002")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var msg = make([]byte, 0, 24)
	// メッセージを送信する
	for {
		time.Sleep(10 * time.Millisecond)
		// Stop
		select {
		case <-stopCh:
			log.Println("(goroutine sendAttitude) stop request received")
			return
		default:
			rollB := float64ToBytes(*roll)
			pitchB := float64ToBytes(*pitch)
			yawB := float64ToBytes(*yaw)
			msg = []byte{}
			msg = append(msg, rollB...)
			msg = append(msg, pitchB...)
			msg = append(msg, yawB...)
			// log.Println(bytesToFloat64(msg[:8]))
			if err := conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
				log.Println(err)
			}
			if _, err := conn.Write(msg); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
