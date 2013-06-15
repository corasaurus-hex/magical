package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println(getHardwareAddressAsUint64())
}

func logFatalErrorf(format string, a ...interface{}) {
	log.Fatal(fmt.Errorf(format, a...))
}
