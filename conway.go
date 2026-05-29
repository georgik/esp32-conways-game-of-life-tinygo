package main

import (
	"machine"
	"math/rand"
	"time"
)

const (
	width  = 30
	height = 30
)

type grid [height][width]int

var g grid

func (g *grid) countCell(x, y int) int {
	var c int
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := x + dx
			ny := y + dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				if g[ny][nx] > 0 {
					c++
				}
			}
		}
	}
	return c
}

func (g *grid) nextStep() grid {
	var ng grid
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := g.countCell(x, y)
			if g[y][x] > 0 {
				if c == 2 || c == 3 {
					ng[y][x] = g[y][x] + 1
				} else {
					ng[y][x] = 0
				}
			} else {
				if c == 3 {
					ng[y][x] = 1
				}
			}
		}
	}
	return ng
}

func (g *grid) randomize() {
	rng, _ := machine.GetRNG()
	rand.Seed(int64(rng))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Intn(100) < 35 {
				g[y][x] = 1
			} else {
				g[y][x] = 0
			}
		}
	}
}

var outBuf [1024]byte
var outIdx int

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

	g.randomize()

	for {
		outIdx = 0
		serialClearScreen()
		serialMoveCursor(1, 1)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if g[y][x] == 0 {
					serialPut(' ')
				} else if g[y][x] <= 9 {
					serialPut(byte('0' + g[y][x]))
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

		g = g.nextStep()
		time.Sleep(250 * time.Millisecond)

		iteration++
		if iteration >= maxIterations {
			g.randomize()
			iteration = 0
		}
	}
}
