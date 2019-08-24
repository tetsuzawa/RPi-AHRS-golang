package server

import (
	"github.com/westphae/quaternion"
	"log"
	"math"
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
	go updateParams(&params)
	go sendAttitude(&roll, &pitch, &yaw)

	// params.updateParams()
	st := time.Now()
	i := 1

	for {
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

/*
func float64ToString(f []float64) []string {
	s := make([]string, len(f))
	for n := range f {
		s[n] = strconv.FormatFloat(f[n], 'f', 8, 64)
	}
	return s
}
*/

/*
func radianToDegreeSlice(rad []float64) (float64, float64, float64) {
	deg := make([]float64, len(rad))
	for n := range rad {
		deg[n] = rad[n] * 180 / math.Pi
	}
	return deg[0], deg[1], deg[2]
}
*/
