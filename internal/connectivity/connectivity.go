package connectivity

import (
	"log"
	"net"
	"os"
)

var (
	connLogger = log.New(os.Stdout, "[ConnCheck]  ", log.LstdFlags)
)

func Check() (bool, bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		connLogger.Fatal(err)
	}
	var v4, v6 bool = false, false
	connLogger.Print("Checking for IPv4 and IPv6 connectivity...")
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		// IPv4
		if ok && ipnet.IP.To4() != nil {
			if !ipnet.IP.IsLoopback() {
				connLogger.Printf("IPv4 connectivity detected: %s", ipnet.IP.String())
				v4 = true
			}
		}
		// IPv6
		if ok && ipnet.IP.To16() != nil {
			if !ipnet.IP.IsLoopback() && !ipnet.IP.IsPrivate() && !ipnet.IP.IsLinkLocalMulticast() && !ipnet.IP.IsLinkLocalUnicast() {
				connLogger.Printf("IPv6 connectivity detected: %s", ipnet.IP.String())
				v6 = true
			}
		}
	}
	return v4, v6
}
