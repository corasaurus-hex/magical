package main

import (
	"net"
	"regexp"
	"strings"
	"testing"
)

var (
	macColon       = "34:B6:02:61:DE:1B"
	macDash        = "34-B6-02-61-DE-1B"
	macDot         = "34B6.0261.DE1B"
	macStripRegexp = regexp.MustCompile("[:.-]")
)

func TestGetHardwareAddr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not get hardware address: %v", r)
		}
	}()

	var hardwareAddr net.HardwareAddr = getHardwareAddr()

	if len(hardwareAddr) == 0 {
		t.Errorf("Got a bad hardware address: %v", hardwareAddr)
	}
}

func TestHardwareAddrAsUint64(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not get hardware address as an integer: %v", r)
		}
	}()

	var getHardwareAddrAsUint64 uint64 = getHardwareAddrAsUint64()

	if getHardwareAddrAsUint64 < 1 {
		t.Errorf("Could not get the hardware address as an integer: %v", getHardwareAddrAsUint64)
	}
}

func TestHardwareAddrToUint64WithGoodValue(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not convert net.HardwareAddr to integer: %v", r)
		}
	}()

	hardwareAddr, _ := net.ParseMAC(macColon)

	uintHardwareAddr := hardwareAddrToUint64(hardwareAddr)

	if uintHardwareAddr != 57956328660507 {
		t.Errorf("Expected %v == %v", "34:B6:02:61:DE:1B", 57956328660507)
	}
}

func BenchmarkMilliseconds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getTimeInMilliseconds()
	}
}

func BenchmarkHardwareAddressAsUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getHardwareAddrAsUint64()
	}
}

func BenchmarkMergeNumbers(b *testing.B) {
	x, y, z := getTimeInMilliseconds(), getHardwareAddrAsUint64(), uint64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mergeNumbers(x, y, z)
	}
}

func BenchmarkRegexpReplaceAll(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		macStripRegexp.ReplaceAllLiteralString(macColon, "")
	}
}

func BenchmarkStringReplaceAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Replace(macColon, ":", "", -1)
		strings.Replace(macDash, ".", "", -1)
		strings.Replace(macDot, "-", "", -1)
	}
}

func BenchmarkStringReplaceSome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Replace(macColon, ":", "", 5)
		strings.Replace(macDash, ".", "", 5)
		strings.Replace(macDot, "-", "", 2)
	}
}
