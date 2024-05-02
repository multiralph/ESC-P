package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Println("Usage: go run main.go <number> <letter> <start> <end> <step>")
		return
	}

	number := os.Args[1]
	letter := os.Args[2]
	start, _ := strconv.Atoi(os.Args[3])
	end, _ := strconv.Atoi(os.Args[4])
	step, _ := strconv.Atoi(os.Args[5])

	printerIP := "at1200-lp-be16" // Replace with the actual IP if necessary
	printerPort := "9100"         // Common port for HP JetDirect / network printers

	conn, err := net.Dial("tcp", printerIP+":"+printerPort)
	if err != nil {
		fmt.Println("Failed to connect to printer:", err)
		return
	}
	defer conn.Close()

	for i := start; i <= end; i += step {
		barcodeText := fmt.Sprintf("%s.%s.%d", number, letter, i)
		escpCommands := []byte{
			0x1B, 0x40, // ESC @ - Initialize ESC/P mode
			0x1B, 0x69, 0x4C, 0x01, // ESC i a 00 - Select ESC/P mode

			0x1B, 0x28, 0x43, 0x02, 0x00, 0x08, 0x07, // Page length setting

			0x1b, 0x69,
			0x72, 0x00, // no text under code
			0x68, 0x20, 0x01, // height
			0x77, 0x03, // width
			0x42, // start content
		}

		for _, c := range barcodeText {
			escpCommands = append(escpCommands, byte(c))
		}

		escpCommands = append(escpCommands, []byte{
			0x5c,             // end barcode
			0x1b, 0x74, 0x02, // character set western europe
			0x1b, 0x6b, 0x09, // Letter Gothic Outline
			0x1B, 0x58, 0x00, 0xE9, 0x00, // character size 44
			0x1b, 0x45, // bold
			0x1B, 0x28, 0x56, 0x02, 0x00, 0xA0, 0x01, // absolut vertial position
			0x1B, 0x24, 0x90, 0x00, // absolut horizontal position
		}...)

		for _, c := range barcodeText {
			escpCommands = append(escpCommands, byte(c))
		}

		escpCommands = append(escpCommands, 0x0C) // FF - Form feed to eject the page

		_, err = conn.Write(escpCommands)
		if err != nil {
			fmt.Println("Failed to send data to printer:", err)
			return
		}
		fmt.Printf("Label sent to printer successfully: %s\n", barcodeText)
	}
}
