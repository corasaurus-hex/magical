package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	timeInMs       uint64
	hardwareAddr   uint64
	sequence       = uint64(0)
	macStripRegexp = regexp.MustCompile(`[^a-fA-F0-9]`)
	mutex          = new(sync.Mutex)
)

func main() {
	timeInMs = getTimeInMilliseconds()
	hardwareAddr = getHardwareAddrUint64()
	for i := 0; i < 1000; i++ {
		fmt.Printf("New value: %v\n", nextId())
	}
}

func getHardwareAddrUint64() uint64 {
	ifs, err := net.Interfaces()

	if err != nil {
		log.Fatalf("Could not get any network interfaces: %v, %+v", err, ifs)
	}

	var hwAddr net.HardwareAddr

	for _, i := range ifs {
		if len(i.HardwareAddr) > 0 {
			hwAddr = i.HardwareAddr
			break
		}
	}

	if hwAddr == nil {
		log.Fatalf("No interface found with a MAC address: %+v", ifs)
	}

	mac := hwAddr.String()
	hex := macStripRegexp.ReplaceAllLiteralString(mac, "")

	u, err := strconv.ParseUint(hex, 16, 64)

	if err != nil {
		log.Fatalf("Unable to parse %v (from mac %v) as an integer: %v", hex, mac, err)
	}

	return u
}

func getTimeInMilliseconds() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}

func mergeNumbers(now uint64, mac uint64, seq uint64) string {
	return fmt.Sprintf("%012x%016x%04x", now, mac, seq)
}

func nextId() string {
	mutex.Lock()
	defer mutex.Unlock()

	newTimeInMs := getTimeInMilliseconds()

	if newTimeInMs == timeInMs {
		sequence += 1
	} else {
		timeInMs = newTimeInMs
		sequence = 0
	}

	return mergeNumbers(timeInMs, hardwareAddr, sequence)
}
