// Package utils provides miscellaneous utilities and cert generation routines
// for Rhine.
package utils

import (
	"log"
	"net"
	"os"
	"path/filepath"
)

// BinDir contains the directory in which executable executed is in.
var BinDir = func() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}()

// Check if an error has occurred and panics if it has.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// GetOutboundIP gets preferred outbound ip of this machine.
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
