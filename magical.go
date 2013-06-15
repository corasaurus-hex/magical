package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"strconv"
	"time"
)

const (
	timestampBits   = uint64(64)
	macAddressBits  = uint64(48)
	sequenceBits    = uint64(16)

	macAddressShift = sequenceBits
	timestampShift  = sequenceBits + macAddressBits
)

func main() {
	fmt.Println(getMilliseconds())
	fmt.Println(getHardwareAddressAsUint64())
}

func logFatalErrorf(format string, a ...interface{}) {
	log.Fatal(fmt.Errorf(format, a...))
}

func getHardwareAddressAsUint64() (uintHardwareAddr uint64) {
	return convertHardwareAddressToUint64(getHardwareAddress())
}

func convertHardwareAddressToUint64(hardwareAddress net.HardwareAddr) (uintHardwareAddr uint64) {
	strippedHardwareAddr := strings.Replace(hardwareAddress.String(), ":", "", -1)

	uintHardwareAddr, err := strconv.ParseUint(strippedHardwareAddr, 16, 48)

	if err != nil {
		logFatalErrorf("Unable to parse %q as an integer: %q", strippedHardwareAddr, err)
	}

	return uintHardwareAddr
}

func getHardwareAddress() (net.HardwareAddr) {
	interfaces, err := net.Interfaces()

	if err != nil {
		logFatalErrorf("Could not get any network interfaces", err)
	}

	for _, value := range interfaces {
		if len(value.HardwareAddr) > 0 {
			return value.HardwareAddr
		}
	}

	logFatalErrorf("No interface found with a MAC address: %v", interfaces)

	return nil
}

func getMilliseconds() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}
