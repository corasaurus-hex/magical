package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"
	"math/big"
)

const (
	macAddressBits = uint(48)
	sequenceBits   = uint(16)
)

var (
	macStripRegexp = regexp.MustCompile("[:.-]")
)

func main() {
	fmt.Println(nextId())
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

func hardwareAddr() net.HardwareAddr {
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

func sequence() uint64 {
	return uint64(0)
}

func mergeNumbers(now uint64, mac uint64, seq uint64) string {
	i := new(big.Int)
	i.SetUint64(now)
	i.Lsh(i, macAddressBits)
	i.Or(new(big.Int).SetUint64(mac), i)
	i.Lsh(i, sequenceBits)
	i.Or(new(big.Int).SetUint64(seq), i)
	return i.String()
}

func nextId() string {
	return mergeNumbers(milliseconds(), hardwareAddrAsUint64(), sequence())
}