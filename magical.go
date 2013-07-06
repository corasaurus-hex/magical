package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maxIds = 10
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

	http.HandleFunc("/", serveIds)
	http.ListenAndServe(":8080", nil)
}

func serveIds(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.ParseInt(r.FormValue("count"), 0, 0)

	if count <= 0 || err != nil {
		count = 1
	} else if count > maxIds {
		count = maxIds
	}

	ids, err := generateIds(int(count))

	if err != nil {
		w.WriteHeader(503)
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, strings.Join(ids, "\n"))
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

func generateIds(count int) ([]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	newTimeInMs := getTimeInMilliseconds()

	if newTimeInMs > timeInMs {
		timeInMs = newTimeInMs
		sequence = 0
	} else if newTimeInMs < timeInMs {
		return nil, fmt.Errorf("Time has reversed! Old time: %v - New time: %v", timeInMs, newTimeInMs)
	}

	ids := make([]string, count)

	for i := 0; i < count; i++ {
		sequence++
		ids[i] = mergeNumbers(timeInMs, hardwareAddr, sequence)
	}

	return ids, nil
}
