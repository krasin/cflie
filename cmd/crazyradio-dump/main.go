// This utility dumps the Flash contents of the flie.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/samofly/crazyradio"
	"github.com/samofly/crazyradio/usb"
)

const BootloaderChannel = 110

func main() {
	buf := make([]byte, 32)

	list, err := usb.ListDevices()
	if err != nil {
		log.Fatalf("Unable list Crazyradio dongles: %v", err)
	}

	info := list[0]
	dev, err := usb.Open(info)
	if err != nil {
		log.Fatalf("Unable to open Crazyradio USB dongle: %v", err)
	}

	err = dev.SetRateAndChannel(crazyradio.DATA_RATE_2M, BootloaderChannel)
	if err != nil {
		log.Fatal("SetRateAndChannel: %v", err)
	}
	for {
		_, err = dev.Write([]byte{0xFF, 0xFF, 0x10})
		if err != nil {
			log.Printf("write: %v", err)
			continue
		}
		n, err := dev.Read(buf)
		if err != nil {
			log.Printf("read: n: %d, err: %v", n, err)
			continue
		}
		if n == 1 && buf[0] == 0 {
			// Empty packet, compact log
			fmt.Fprintf(os.Stderr, ".")
			continue
		}
		log.Printf("Packet: %v", buf[:n])
	}
}
