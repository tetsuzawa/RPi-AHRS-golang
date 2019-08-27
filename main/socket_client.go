package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

func float64ToBytes(f float64) []byte {
	bits := math.Float64bits(f)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)
	return b
}
func bytesToFloat64(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	f := math.Float64frombits(bits)
	return f
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, f); err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

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
		case <- stopCh:
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
