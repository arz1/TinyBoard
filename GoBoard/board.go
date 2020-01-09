package main

//Build + Deployment
//tinygo build -opt=1 -o=/media/adam/JLINK/flash.hex -target=pca10056 hello.go
//tinygo build -opt=1 -o=./flash.elf -target=pca10056 board.go
//nrfjprog --reset -f nrf52

import (
	//"fmt"
	//"github.com/enriquebris/goconcurrentqueue"
	//"fmt"
	"machine"
	"time"
)

func main() {

	SetupMachine()

	for {

		frameCapacity := 16
		frame := make([]byte, frameCapacity)
		inserted := 0

		n, buf := UARTReadBytes(frameCapacity + 1)
		if n > 0 {
			for i := range buf {
				if buf[i] != 0x03 {
					if inserted >= frameCapacity {
						//panic("ERROR - Frame capacity exceeded.")
						continue
					}
					frame[inserted] = buf[i]
					inserted++
					//fmt.Println(string(frame))
				} else {
					reversed := reverse(string(frame))
					UARTWriteBytes([]byte(reversed))
					UARTWriteBytes([]byte{0x03})
					//fmt.Println("SENT:  " + string(reversed))
					frame = make([]byte, frameCapacity)
					inserted = 0
				}
			}
		}
	}
}

func SetupMachine() {
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{BaudRate: 9600, TX: 6, RX: 8})
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led4 := machine.LED4
	led4.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.LED.High()
	machine.LED4.High()
}

func UARTWriteBytes(buffer []byte) {

	led := machine.LED
	led.Low()

	machine.UART0.Write(buffer)

	time.Sleep(time.Millisecond * 100)
	led.High()
}

func UARTWrite(buffer string) {

	UARTWriteBytes([]byte(buffer))
}

func UARTReadBytes(bufferSize int) (int, []byte) {

	machine.LED4.Low()

	buf := make([]byte, bufferSize)

	n, err := machine.UART0.Read(buf)

	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 100)
	machine.LED4.High()
	return n, buf
}

func UARTRead(bufferSize int) (int, string) {

	n, buffer := UARTReadBytes(bufferSize)
	text := string(buffer)
	return n, text
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
