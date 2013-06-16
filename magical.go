package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
	"regexp"
)

const (
	timestampBits   = uint64(64)
	macAddressBits  = uint64(48)
	sequenceBits    = uint64(16)

	macAddressShift = sequenceBits
	timestampShift  = sequenceBits + macAddressBits
)

var (
	macStripRegexp = regexp.MustCompile("[:.-]")
)

func main() {
	fmt.Println(milliseconds())
	fmt.Println(hardwareAddrAsUint64())
}

func hardwareAddrAsUint64() (uintHardwareAddr uint64) {
	return convertHardwareAddrToUint64(hardwareAddr())
}

func convertHardwareAddrToUint64(hardwareAddr net.HardwareAddr) (uintHardwareAddr uint64) {
	strippedHardwareAddr := macStripRegexp.ReplaceAllLiteralString(hardwareAddr.String(), "")

	uintHardwareAddr, err := strconv.ParseUint(strippedHardwareAddr, 16, 48)

	if err != nil {
		log.Fatalf("Unable to parse %q as an integer: %q", strippedHardwareAddr, err)
	}

	return uintHardwareAddr
}

func hardwareAddr() (net.HardwareAddr) {
	interfaces, err := net.Interfaces()

	if err != nil {
		log.Fatalf("Could not get any network interfaces", err)
	}

	for _, value := range interfaces {
		if len(value.HardwareAddr) > 0 {
			return value.HardwareAddr
		}
	}

	log.Fatalf("No interface found with a MAC address: %v", interfaces)

	return nil
}

func milliseconds() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}
