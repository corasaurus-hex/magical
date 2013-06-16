package main

import (
	"testing"
	"net"
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
