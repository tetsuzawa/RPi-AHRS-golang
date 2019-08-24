package server

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
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
	// f := make([]float64, len(ssep))
	f := stringToFloat64(ssep)
	if len(f) != 9 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, errors.New("received data is incorrect")
	}
	return f[0], f[1], f[2],
		f[3], f[4], f[5], f[6], f[7], f[8], nil
}

func updateParams(p *Parameters) {

	conn, err := net.ListenPacket("udp", "127.0.0.1:62000")
	if err != nil {
		println("ERROR: ", err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	var data string

	for {
		time.Sleep(5 * time.Millisecond)
		length, remoteAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatalln("Connect ERROR : ", err)
		}
		data = string(buffer[:length])
		log.Printf("\nReceived from %v: %v\n", remoteAddr, data)
		/*
			if len(data) == 0 {
				continue
			}
		*/
		p.wx, p.wy, p.wz, p.ax, p.ay, p.az, p.mx, p.my, p.mz, err = splitDataToFloat64(data, ",")
		if err != nil {
			log.Fatalln("Value Error:", err)
		}
	}
}

/*
func (p *Parameters) updateParams() {
	conn, err := net.ListenPacket("udp", "127.0.0.1:62001")
	if err != nil {
		println("ERROR: ", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	data := "0.01,0.01,0.01,0.01,0.01,1.00,30.0,0,0"

	// 通信読込 + 接続相手アドレス情報が受取
	go func(data *string) {
		for {
			time.Sleep(20 * time.Millisecond)
			length, remoteAddr, _ := conn.ReadFrom(buffer)
			log.Println(len(*data))
			*data = string(buffer[:length])
			log.Println(len(*data))
			if len(*data) == 0 {
				continue
			}
			p.wx, p.wy, p.wz, p.ax, p.ay, p.az, p.mx, p.my, p.mz, err = splitDataToFloat64(*data, ",")
			if err != nil {
				continue
			}
			log.Printf("Received from %v: %v\n", remoteAddr, data)
		}
	}(&data)

}
*

/*
func (p *Parameters) updateParams() {
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	conn, err := net.ListenPacket("udp", "127.0.0.1:62000")
	if err != nil {
		println("ERROR: ", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	data := "initial"

	// 通信読込 + 接続相手アドレス情報が受取
	go func(stopCh, doneCh chan struct{}) {
		defer func() { close(doneCh) }()
		for {
			length, remoteAddr, _ := conn.ReadFrom(buffer)
			data = string(buffer[:length])
			p.wx, p.wy, p.wz, p.ax, p.ay, p.az, p.mx, p.my, p.mz = splitDataToFloat64(data, ",")
			log.Printf("Received from %v: %v\n", remoteAddr, data)
			select {
			case <-stopCh:
				println("stop request received.")
				return
			default:
			}
		}
	}(stopCh, doneCh)

	// time.Sleep(30 * time.Second)
	println("request stop.")
	close(stopCh)

	// loopが完了するまで待つ
	<-doneCh
	println("loop done.")
￿}
*/

/*
func main_socket() {
	fmt.Println("Server is Running at 127.0.0.1:62000")
	conn, err := net.ListenPacket("udp", "127.0.0.1:62000")
	if err != nil{
		println("ERROR: ", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	data := "initial"
	i:=1

	// 通信読込 + 接続相手アドレス情報が受取
	go func(data *string) {
		for {
			length, remoteAddr, _ := conn.ReadFrom(buffer)
			*data = string(buffer[:length])
			fmt.Printf("Received from %v: %v\n", remoteAddr, *data)
			conn.WriteTo([]byte("Hello, World !"), remoteAddr)
		}
	}(&data)

	for{
		println(i, data)
		time.Sleep(20 * time.Millisecond)
		i++
	}

}

*/
