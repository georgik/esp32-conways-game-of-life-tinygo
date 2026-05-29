# Conway's Game of Life - ESP32-S3 (TinyGo)

Conway's Game of Life implemented in Go for the ESP32-S3 microcontroller using TinyGo. The simulation is streamed over UART and displayed in a terminal as ANSI escape sequences. Live cells are rendered with their age (1-9 or `#` for age 10+).

## Hardware Requirements

- ESP32-S3 development board
- USB-C cable for programming and serial communication

## Software Requirements

- [TinyGo](https://tinygo.org/)
- A terminal emulator that supports ANSI escape sequences (e.g., `screen`, `minicom`, `picocom`, or macOS Terminal.app)

## Building

Compile the firmware and flash it to the device:

```
tinygo flash -target esp32s3-generic conway.go
```

This compiles the program and automatically uploads it to the connected ESP32-S3 board.

## Serial Monitor

To view the simulation output, open a serial monitor:

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

## Development

Before committing, ensure the code follows Go formatting standards:

```
go fmt conaw.go
```

