# Conway's Game of Life - ESP32-S3 (TinyGo)

Conway's Game of Life implemented in Go for ESP32-S3 microcontrollers using TinyGo. Multiple board implementations available with different display options.

## Project Structure

```
conway/              # Shared game logic package
├── grid.go
boards/
├── esp32s3-generic/ # UART output (terminal display)
│   └── main.go
└── esp32s3-box3/    # ILI9341 LCD display
    └── main.go
```

## Hardware Requirements

### UART Version (esp32s3-generic)
- Any ESP32-S3 development board
- USB-C cable for programming and serial communication

### Display Version (esp32s3-box3)
- ESP32-S3-BOX-3 or compatible board with ILI9341 display (320x240)
- USB-C cable

## Software Requirements

- [TinyGo](https://tinygo.org/) 0.42.0+
- For UART: terminal emulator with ANSI support (`screen`, `minicom`, `picocom`)

## Building

### UART Version (Terminal Output)

**Compile and flash:**
```
tinygo flash -target esp32s3-generic ./boards/esp32s3-generic/main.go
```

**Build only:**
```
tinygo build -target esp32s3-generic -o conway-uart.bin ./boards/esp32s3-generic/main.go
```

### Display Version (ILI9341 LCD)

**Note:** ESP32-S3-BOX-3 requires manual entry into download mode:
1. Hold BOOT button
2. Press RESET button
3. Release both buttons
4. Run flash command

**Compile and flash:**
```
tinygo flash -target esp32s3-generic ./boards/esp32s3-box3/main.go
```

**Build only:**
```
tinygo build -target esp32s3-generic -o conway-display.bin ./boards/esp32s3-box3/main.go
```

For boards with dedicated `esp32s3-box3` target:
```
tinygo flash -target esp32s3-box3 ./boards/esp32s3-box3/main.go
```

For debugging, build ELF with symbols:
```
tinygo build -target esp32s3-generic -o conway.elf ./boards/esp32s3-generic/main.go
```

## Serial Monitor (UART Version)

To view the UART simulation output:

```
tinygo monitor
```

Or using `screen`:

```
screen /dev/cu.usbmodem1401 115200
```

Replace the serial device path as needed for your system. On Linux, the port is typically `/dev/ttyACM0` or `/dev/ttyUSB0`.

## How It Works

The program implements Conway's Game of Life on a 30x30 grid. Each live cell tracks its age:

- **` ` (space)** -- Dead cell
- **`1` - `9`** -- Live cell, age equals the number (how many generations it has survived)
- **`#`** -- Live cell older than 9 generations

A new simulation starts automatically after 200 iterations. The grid is initialized with a random pattern at approximately 35% cell density.

### Display Rendering

**UART Version:** ANSI escape sequences render colored numbers (1-9) and `#` for age 10+.

**Display Version:** Color gradient on ILI9341 LCD:
- Dead: Black
- Young (1-2): Dark green
- Medium (3-4): Green
- Mature (5-6): Yellow-green
- Old (7-8): Yellow
- Ancient (9+): Red-purple

## Development

Format code before committing:

```
go fmt ./conway/ ./boards/
```

