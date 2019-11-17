package main

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func radianToDegree(rad float64) (deg float64) {
	deg = rad * 180 / math.Pi
	return
}

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
