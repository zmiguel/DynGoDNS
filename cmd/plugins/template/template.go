package main

import (
	"log"
	"os"

	"github.com/zmiguel/DynGoDNS/internal/types"
)

var (
	//api_endpoint   = "https://api.cloudflare.com/client/v4"
	Config         types.Config
	templateLogger = log.New(os.Stdout, "[Template]   ", log.LstdFlags)
) //     always leave this much space --> "             "  thank you.

func Info() string {
	return "Template  v0.0.1"
}

func Initialise() {
	//necessary init stuff here, called once at startup
}

func GetV4(dom string) (bool, string, string) {
	// get current DNS records
	templateLogger.Print("Getting current DNS records")
	return false, "", ""
}

func CreateV4(dom string, ip string) {
	// Create new DNS record
	templateLogger.Printf("Creating new DNS record for %s with ip %s", dom, ip)
}

func UpdateV4(dom string, ip string, id string) {
	// Update DNS record
	templateLogger.Printf("Updating DNS record for %s with ip %s", dom, ip)
}

func DeleteV4(dom string, ip string, id string) {
	// Update DNS record
	templateLogger.Printf("Updating DNS record for %s with ip %s", dom, ip)
}

func GetV6(dom string) (bool, string, string) {
	// get current DNS records
	templateLogger.Print("Getting current DNS records")
	return false, "", ""
}

func CreateV6(dom string, ip string) {
	// Create new DNS record
	templateLogger.Printf("Creating new DNS record for %s with ip %s", dom, ip)
}

func UpdateV6(dom string, ip string, id string) {
	// Update DNS record
	templateLogger.Printf("Updating DNS record for %s with ip %s", dom, ip)
}

func DeleteV6(dom string, ip string, id string) {
	// Update DNS record
	templateLogger.Printf("Updating DNS record for %s with ip %s", dom, ip)
}
