package main

import (
	"testing"
	"net"
)

func TestHardwareAddress(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not get hardware address: %v", r)
		}
	}()

	var hardwareAddress net.HardwareAddr = hardwareAddress()

	if len(hardwareAddress) == 0 {
		t.Errorf("Got a bad hardware address: %v", hardwareAddress)
	}
}

func TestHardwareAddressAsUint64(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not get hardware address as an integer: %v", r)
		}
	}()

	var hardwareAddressAsUint64 uint64 = hardwareAddressAsUint64()

	if hardwareAddressAsUint64 < 1 {
		t.Errorf("Could not get the hardware address as an integer: %v", hardwareAddressAsUint64)
	}
}

func TestConvertHardwareAddressToUint64WithGoodValue(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Could not convert net.HardwareAddr to integer: %v", r)
		}
	}()

	hardwareAddr, _ := net.ParseMAC("34:B6:02:61:DE:1B")

	uintHardwareAddr := convertHardwareAddressToUint64(hardwareAddr)

	if uintHardwareAddr != 57956328660507 {
		t.Errorf("Expected %v == %v", "34:B6:02:61:DE:1B", 57956328660507)
	}
}
