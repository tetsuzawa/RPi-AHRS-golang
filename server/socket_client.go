package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
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

func sendAttitude(roll, pitch, yaw *float64) {
	conn, err := net.Dial("udp", "127.0.0.1:60002")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var msg = make([]byte, 0, 24)
	// メッセージを送信する
	for {
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
