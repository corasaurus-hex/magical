package main

import (
	"net"
	"strings"
	"strconv"
)

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
