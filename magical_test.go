package main

import (
	"strings"
	"testing"
)

var (
	macColon = "34:B6:02:61:DE:1B"
	macDash  = "34-B6-02-61-DE-1B"
	macDot   = "34B6.0261.DE1B"
)

func setup() {
	timeInMs = getTimeInMilliseconds()
	hardwareAddr = getHardwareAddrUint64()
}

func TestGetHardwareAddr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not get hardware address: %v", r)
		}
	}()

	getHardwareAddrUint64()
}

func BenchmarkMilliseconds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getTimeInMilliseconds()
	}
}

func BenchmarkHardwareAddressAsUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getHardwareAddrUint64()
	}
}

func BenchmarkMergeNumbers(b *testing.B) {
	x, y, z := getTimeInMilliseconds(), getHardwareAddrUint64(), uint64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mergeNumbers(x, y, z)
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

func BenchmarkNextId(b *testing.B) {
	setup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nextId()
	}
}