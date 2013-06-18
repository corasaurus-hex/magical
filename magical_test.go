package main

import (
	"net"
	"testing"
)

func TestHardwareAddr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not get hardware address: %v", r)
		}
	}()

	var hardwareAddr net.HardwareAddr = hardwareAddr()

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

	var hardwareAddrAsUint64 uint64 = hardwareAddrAsUint64()

	if hardwareAddrAsUint64 < 1 {
		t.Errorf("Could not get the hardware address as an integer: %v", hardwareAddrAsUint64)
	}
}

func TestConvertHardwareAddrToUint64WithGoodValue(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not convert net.HardwareAddr to integer: %v", r)
		}
	}()

	hardwareAddr, _ := net.ParseMAC("34:B6:02:61:DE:1B")

	uintHardwareAddr := convertHardwareAddrToUint64(hardwareAddr)

	if uintHardwareAddr != 57956328660507 {
		t.Errorf("Expected %v == %v", "34:B6:02:61:DE:1B", 57956328660507)
	}
}

func BenchmarkMilliseconds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		milliseconds()
	}
}

func BenchmarkHardwareAddressAsUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hardwareAddrAsUint64()
	}
}

func BenchmarkSequence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sequence()
	}
}

func BenchmarkMergeNumbers(b *testing.B) {
	x, y, z := milliseconds(), hardwareAddrAsUint64(), sequence()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mergeNumbers(x, y, z)
	}
}
