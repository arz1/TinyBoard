package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)

var sPort *serial.Port = nil

func main() {

	SetupPort("/dev/ttyACM0")

	for {
		//	fmt.Println("\nPress the Enter Key to continue...")
		//	var input string
		//	fmt.Scanln(&input)

		dataToSent := "0123456789abcdef"
		fmt.Printf("\nData to sent...  " + dataToSent)
		fmt.Printf("\nHEX: " + TextToHexText(dataToSent))

		WriteBytes([]byte(dataToSent))
		WriteBytes([]byte{0x03})

		fmt.Printf("\nWaiting for response & Reading...\n")
		time.Sleep(time.Millisecond * 100)
		//for {
		//	_, r1 := Read(100)
		//	fmt.Printf(string(r1))
		//}

		frameCapacity := 16
		frame := make([]byte, frameCapacity)
		inserted := 0

		n, buf := ReadBytes(frameCapacity + 1)
		if n > 0 {

			mystr := string(buf)
			if mystr == "s" {
				panic("")
			}
			for i := range buf {
				if buf[i] != 0x03 {
					if inserted >= frameCapacity {
						//panic("PC: ERROR - Frame capacity exceeded.")
						continue
					}
					frame[inserted] = buf[i]
					inserted++
					//fmt.Println(string(frame))
				} else {

					fmt.Printf("\nResponse from board: " + string(frame))
					fmt.Printf("\nHEX: " + TextToHexText(string(frame)))
					fmt.Printf("\n==================================================")

					//fmt.Println("SENT:  " + string(reversed))
					frame = make([]byte, frameCapacity)
					inserted = 0
				}
			}
		}

		time.Sleep(time.Millisecond * 800)
	}
}

func SetupPort(portId string) {

	portConfig := &serial.Config{Name: portId, Baud: 9600}
	var err error = nil
	sPort, err = serial.OpenPort(portConfig)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func TextToHexText(plainText string) string {
	hexText := fmt.Sprintf("% x\n", plainText)
	return hexText
}

func WriteBytes(buffer []byte) int {

	n, err := sPort.Write(buffer)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return n
}

func Write(buffer string) int {

	n := WriteBytes([]byte(buffer))
	return n
}

func ReadBytes(bufferSize int) (int, []byte) {

	buf := make([]byte, bufferSize)

	n, err := sPort.Read(buf)

	if err != nil {
		panic(err)
	}
	return n, buf
}

func Read(bufferSize int) (int, string) {

	n, buffer := ReadBytes(bufferSize)
	text := string(buffer)
	return n, text
}
