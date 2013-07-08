package main

import (
	// "fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

var (
	macColon    = "34:B6:02:61:DE:1B"
	macDash     = "34-B6-02-61-DE-1B"
	macDot      = "34B6.0261.DE1B"
	hexIdRegexp = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
)

func TestGetHardwareAddr(t *testing.T) {
	setup()
	hw := getHardwareAddrUint64()
	kind := reflect.TypeOf(hw).Kind()

	if kind != reflect.Uint64 {
		t.Fatalf("Expected %v to equal %v", kind, reflect.Uint64)
	}

	if hw <= 0 {
		t.Fatalf("Expected %v to be greater than 0", hw)
	}
}

func TestGetTimeInMilliseconds(t *testing.T) {
	setup()
	now := getTimeInMilliseconds()
	time.Sleep(1 * time.Millisecond)
	later := getTimeInMilliseconds()

	difference := later - now

	// sleep + time to get next milliseconds may be > 1ms
	if difference > 2 || difference < 1 {
		t.Fatalf("Expected %v to be 1 or 2 larger than %v", later, now)
	}
}

func TestGenerateIds(t *testing.T) {
	setup()
	idGeneration := map[int]int {
		-2: 1,
		-1: 1,
		0: 1,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
		10: 10,
		11: 10,
		12: 10,
		13: 10,
	}
	for arg, length := range idGeneration {
		generatedIds, err := generateIds(arg)

		if err != nil {
			t.Fatalf("Expected generateIds(%v) to generate %v ids, instead there was an error: %v", arg, length, err)
		}

		if len(generatedIds) != length {
			t.Fatalf("Expected generateIds(%v) to generate %v ids, instead it generated %v ids", arg, length, len(generatedIds))
		}
	}
}

func TestGenerateIdsConcurrently(t *testing.T) {
	setup()

	idsChannel := make(chan []id, 10)

	// startup goroutines that generate a crap-ton of ids
	for i := 0; i < 10; i++ {
		go func() {
			generatedIds := []id{}
			for count := -10; count < 20; count++ {
				ids, err := generateIds(count)
				if err != nil {
					t.Fatalf("Expected generatedIds(%v) to work, got error: %v", count, err)
				}
				generatedIds = append(generatedIds, ids...)
			}
			idsChannel <- generatedIds
		}()
	}

	// gather together all the generated ids
	allGeneratedIds := []string{}

	for i := 0; i < 10; i++ {
		generatedIds := <- idsChannel
		for _, id := range generatedIds {
			allGeneratedIds = append(allGeneratedIds, id.Hex())
		}
	}

	// make sure none are blank
	for i, id := range allGeneratedIds {
		if id == "" {
			t.Fatalf("Blank id found at index: %v", i)
		}
	}

	// no idea how to test uniqueness of a slice, this is my feeble attempt
	sort.Strings(allGeneratedIds)
	lastId := ""

	for i, id := range allGeneratedIds {
		if lastId == id {
			t.Fatalf("Duplicate id found at index %v: %v", i, id)
		}
	}
}

func TestGenerateHexIds(t *testing.T) {
	setup()

	ids, err := generateHexIds(1)

	if err != nil {
		t.Fatalf("Expected generatedIds(1) to work, got error: %v", err)
	}

	if !hexIdRegexp.MatchString(ids[0]) {
		t.Fatalf("generateHexIds(1) did not generate a valid hex id: %v", ids[0])
	}
}

func BenchmarkMilliseconds(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getTimeInMilliseconds()
	}
}

func BenchmarkHardwareAddressAsUint64(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getHardwareAddrUint64()
	}
}

func BenchmarkGeneratingHex(b *testing.B) {
	setup()
	t := id{getTimeInMilliseconds(), getHardwareAddrUint64(), uint64(0)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Hex()
	}
}

func BenchmarkRegexpReplaceAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		macStripRegexp.ReplaceAllLiteralString(macColon, "")
	}
}

func BenchmarkStringReplaceAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Replace(macColon, ":", "", -1)
		strings.Replace(macDash, "-", "", -1)
		strings.Replace(macDot, ".", "", -1)
	}
}

func BenchmarkStringReplaceSome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Replace(macColon, ":", "", 5)
		strings.Replace(macDash, "-", "", 5)
		strings.Replace(macDot, ".", "", 2)
	}
}

func BenchmarkGenerateIds1(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateIds(1)
	}
}

func BenchmarkGenerateIds2(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateIds(2)
	}
}

func BenchmarkGenerateIds5(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateIds(5)
	}
}

func BenchmarkGenerateIds10(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateIds(10)
	}
}

func BenchmarkGenerateHexIds1(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateHexIds(1)
	}
}

func BenchmarkGenerateHexIds2(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateHexIds(2)
	}
}

func BenchmarkGenerateHexIds5(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateHexIds(5)
	}
}

func BenchmarkGenerateHexIds10(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateHexIds(10)
	}
}
