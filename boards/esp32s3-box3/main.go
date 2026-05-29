// Build with: tinygo build -target esp32s3-generic -o conway-display.bin ./boards/esp32s3-box3/main.go
// Or for esp32s3-box3 target when available:
// tinygo build -target esp32s3-box3 -o conway-display.bin ./boards/esp32s3-box3/main.go

// Build with: tinygo build -target esp32s3-generic -o conway-display.bin ./boards/esp32s3-box3/main.go
// Or for esp32s3-box3 target when available:
// tinygo build -target esp32s3-box3 -o conway-display.bin ./boards/esp32s3-box3/main.go

package main

import (
	"image/color"
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/ili9341"
)

const (
	width    = 32
	height   = 24
	cellSize = 10
)

type Grid [height][width]int

func (g *Grid) countCell(x, y int) int {
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

func (g *Grid) nextStep() Grid {
	var ng Grid
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

func (g *Grid) randomize() {
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

var (
	display *ili9341.Device
	g       Grid
)

func getDisplayColor(age int) color.RGBA {
	switch {
	case age == 0:
		return color.RGBA{0, 0, 0, 255}
	case age <= 2:
		return color.RGBA{0, 100, 0, 255}
	case age <= 4:
		return color.RGBA{0, 180, 0, 255}
	case age <= 6:
		return color.RGBA{100, 200, 0, 255}
	case age <= 8:
		return color.RGBA{200, 200, 0, 255}
	default:
		return color.RGBA{200, 0, 100, 255}
	}
}

func initDisplay() {
	const (
		LCD_SCK = machine.GPIO7
		LCD_SDO = machine.GPIO6
		LCD_DC  = machine.GPIO4
		LCD_CS  = machine.GPIO5
		LCD_RST = machine.GPIO48
		LCD_BL  = machine.GPIO47
	)

	rst := machine.GPIO48
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rst.High()
	time.Sleep(10 * time.Millisecond)
	rst.Low()
	time.Sleep(10 * time.Millisecond)

	machine.SPI1.Configure(machine.SPIConfig{
		SCK:       LCD_SCK,
		SDO:       LCD_SDO,
		SDI:       machine.NoPin,
		Frequency: 40e6,
	})

	display = ili9341.NewSPI(
		machine.SPI1,
		LCD_DC,
		LCD_CS,
		machine.NoPin,
	)

	display.Configure(ili9341.Config{
		Width:            320,
		Height:           240,
		DisplayInversion: true,
	})

	display.SetRotation(ili9341.Rotation0Mirror)

	bl := machine.GPIO47
	bl.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bl.High()

	display.FillScreen(color.RGBA{0, 0, 0, 255})
}

func fillCell(x, y int, c color.RGBA) {
	var cellData [100]uint16
	c565 := uint16(c.R>>3)<<11 | uint16(c.G>>2)<<5 | uint16(c.B>>3)
	for i := range cellData {
		cellData[i] = c565
	}
	display.DrawRGBBitmap(int16(x*cellSize), int16(y*cellSize), cellData[:], cellSize, cellSize)
}

func renderGrid() {
	display.FillScreen(color.RGBA{0, 0, 0, 255})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if g[y][x] > 0 {
				fillCell(x, y, getDisplayColor(g[y][x]))
			}
		}
	}
}

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	machine.Serial.Write([]byte("Conway's Game of Life - ESP32-S3-BOX-3\r\n"))
	machine.Serial.Write([]byte("Initializing display...\r\n"))

	initDisplay()

	machine.Serial.Write([]byte("Display initialized!\r\n"))
	machine.Serial.Write([]byte("Starting simulation...\r\n"))

	var iteration int
	const maxIterations = 200

	rng, _ := machine.GetRNG()
	rand.Seed(int64(rng))
	g.randomize()

	for {
		renderGrid()

		g = g.nextStep()
		time.Sleep(100 * time.Millisecond)

		iteration++
		if iteration >= maxIterations {
			rng, _ = machine.GetRNG()
			rand.Seed(int64(rng))
			g.randomize()
			iteration = 0
		}
	}
}
