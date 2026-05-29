package main

import (
	"esp32-conways-game-of-life-tinygo/conway"
	"machine"
	"math/rand"
	"time"
)

var outBuf [1024]byte
var outIdx int
var g conway.Grid

func serialPut(b byte) {
	outBuf[outIdx] = b
	outIdx++
}

func serialFlush() {
	if outIdx > 0 {
		machine.Serial.Write(outBuf[:outIdx])
		outIdx = 0
	}
}

func serialWriteANSIEscape(code string) {
	serialPut(0x1b)
	serialPlace(code)
	serialPut('m')
}

func serialPlace(s string) {
	for i := 0; i < len(s); i++ {
		serialPut(s[i])
	}
}

func serialMoveCursor(x, y int) {
	serialPut(0x1b)
	serialPlace("[")
	if y > 0 {
		serialWriteDigit(y)
	}
	serialPlace(";")
	if x > 0 {
		serialWriteDigit(x)
	}
	serialPlace("H")
}

func serialClearScreen() {
	serialPut(0x1b)
	serialPlace("[2J")
}

func serialWriteDigit(n int) {
	if n >= 100 {
		serialPut(byte('0' + n/100))
		n %= 100
	}
	if n >= 10 {
		serialPut(byte('0' + n/10))
		n %= 10
	}
	serialPut(byte('0' + n))
}

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	time.Sleep(500 * time.Millisecond)

	var iteration int
	const maxIterations = 200

	rng, _ := machine.GetRNG()
	rand.Seed(int64(rng))
	g.RandomizeWith(func() uint32 { return uint32(rand.Uint32()) })

	for {
		outIdx = 0
		serialClearScreen()
		serialMoveCursor(1, 1)

		for y := 0; y < conway.Height; y++ {
			for x := 0; x < conway.Width; x++ {
				cell := g.GetCell(x, y)
				if cell == 0 {
					serialPut(' ')
				} else if cell <= 9 {
					serialPut(byte('0' + cell))
				} else {
					serialPut('#')
				}
			}
			serialPlace("\r\n")
		}

		serialPlace("Iteration: ")
		serialWriteDigit(iteration)
		serialPlace("\r\n")

		serialFlush()

		g = g.NextStep()
		time.Sleep(250 * time.Millisecond)

		iteration++
		if iteration >= maxIterations {
			rng, _ := machine.GetRNG()
			rand.Seed(int64(rng))
			g.RandomizeWith(func() uint32 { return uint32(rand.Uint32()) })
			iteration = 0
		}
	}
}
