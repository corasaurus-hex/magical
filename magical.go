package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"math/big"
)

const (
	macAddressBits = uint(48)
	sequenceBits   = uint(16)
)

func main() {
	fmt.Println(nextId())
}

func getHardwareAddrAsUint64() (uintHardwareAddr uint64) {
	return hardwareAddrToUint64(getHardwareAddr())
}

func hardwareAddrToUint64(h net.HardwareAddr) (uintHardwareAddr uint64) {
	s := h.String()
	s = strings.Replace(s, ":", "", -1)
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, "-", "", -1)

	u, err := strconv.ParseUint(s, 16, 48)

	if err != nil {
		log.Fatalf("Unable to parse %q as an integer: %q", s, err)
	}

	return u
}

func getHardwareAddr() net.HardwareAddr {
	ifs, err := net.Interfaces()

	if err != nil {
		log.Fatalf("Could not get any network interfaces: %v, %+v", err, ifs)
	}

	for _, i := range ifs {
		if len(i.HardwareAddr) > 0 {
			return i.HardwareAddr
		}
	}

	log.Fatalf("No interface found with a MAC address: %+v", ifs)

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
	return mergeNumbers(milliseconds(), getHardwareAddrAsUint64(), sequence())
}
