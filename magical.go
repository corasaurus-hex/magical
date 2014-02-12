package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type idGenerator struct {
	time uint64
	mac  uint64
	seq  *uint64
}

func (i *idGenerator) NextHex() (string, error) {
	seq := atomic.AddUint64(i.seq, 1)

	if seq > maxSequence {
		return "", errors.New("Ran out of numbers for this timestamp")
	}

	t := make([]byte, 8)
	s := make([]byte, 8)
	a := make([]byte, 16)
	binary.BigEndian.PutUint64(t, i.time)
	binary.BigEndian.PutUint64(a[6:14], i.mac)
	binary.BigEndian.PutUint64(s, seq)

	copy(a[0:6], t[2:8])
	copy(a[14:16], s[6:8])

	return hex.EncodeToString(a), nil
}

type idGenerators struct {
	time      uint64
	mac       uint64
	generator *idGenerator
	mutex     *sync.Mutex
}

func (i *idGenerators) GetGenerator(time uint64) (*idGenerator, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if time == i.time {
		return i.generator, nil
	} else if time > i.time {
		i.generator = &idGenerator{time, i.mac, new(uint64)}
		i.time      = time
		return i.generator, nil
	}

	return nil, errors.New("time went backwards")
}

type idGeneratorsPartitions struct {
	generators [100]*idGenerators
}

func (i *idGeneratorsPartitions) Setup() {
	for x := 0; x < 100; x++ {
		igs := &idGenerators{}
		igs.time = getTimeInMilliseconds()
		igs.mac = getHardwareAddrUint64()
		igs.generator = &idGenerator{igs.time, igs.mac, new(uint64)}
		igs.mutex = new(sync.Mutex)
		i.generators[x] = igs
	}
}

func (i *idGeneratorsPartitions) GetGenerator(time uint64) (*idGenerator, error) {
	ig, err := i.generators[time % 100].GetGenerator(time)
	return ig, err
}

func getTimeInMilliseconds() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
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

const (
	defaultCount = 1
	maxIds       = 65535
	maxSequence  = 65535
)

var (
	generators     = &idGeneratorsPartitions{}
	macStripRegexp = regexp.MustCompile(`[^a-fA-F0-9]`)
)

func main() {
	setup()

	http.HandleFunc("/", serveIds)
	http.ListenAndServe(":8080", nil)
}

func setup() {
	generators.Setup()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func serveIds(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.ParseInt(r.FormValue("count"), 0, 0)
	ids, err := generateIds(int(count))

	if err != nil {
		w.WriteHeader(503)
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, strings.Join(ids, "\n"))
}

func generateIds(count int) ([]string, error) {
	if count < 1 {
		count = defaultCount
	} else if count > maxIds {
		count = maxIds
	}

	generator, err := generators.GetGenerator(getTimeInMilliseconds())
	if err != nil {
		return nil, err
	}

	ids := make([]string, count)

	for i := 0; i < count; i++ {
		ids[i], err = generator.NextHex()
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}
