package updater

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/zmiguel/DynGoDNS/internal/connectivity"
	"github.com/zmiguel/DynGoDNS/internal/types"
)

var (
	Config    types.Config
	init_done = false
	dns       types.DNS
	upLogger  = log.New(os.Stdout, "[Updater]    ", log.LstdFlags)
)

func Initialise(d types.DNS, c types.Config) {
	dns = d
	Config = c
}

func Update() {
	if !init_done {
		dns.Initialise.(func())()
		init_done = true
	}

	// Check for IPv4 and IPv6 connectivity
	v4, v6 := connectivity.Check()

	for _, domain := range Config.Domains {
		// get only first part of domain separated by commas
		if strings.Contains(domain, ",") {
			domain = strings.Split(domain, ",")[0]
		}
		upLogger.Print("----------")
		upLogger.Printf("Checking domain: %s ...", domain)
		if Config.V4.Enabled {
			if !v4 {
				upLogger.Print("No IPv4 connectivity, skipping...")
				continue
			}
			//Check current IP
			currentIP := getCurrentIP(4)
			if currentIP == "" {
				upLogger.Printf("No IPv4 found, skipping...")
				continue
			}
			// Get DNS IP
			force, dnsIP, entryID := dns.GetV4.(func(string) (bool, string, string))(domain)
			//  Update DNS IP if needed
			if force {
				upLogger.Printf("Forcing update of DNS record for %s (%s)...", domain, entryID)
				dns.UpdateV4.(func(string, string, string))(domain, currentIP, entryID)
			} else {
				if dnsIP == "" {
					// Create DNS
					upLogger.Printf("No DNS record found for %s, creating...", domain)
					dns.CreateV4.(func(string, string))(domain, currentIP)
				} else if currentIP != dnsIP {
					// Update DNS
					upLogger.Printf("Updating DNS record for %s (%s)...", domain, entryID)
					dns.UpdateV4.(func(string, string, string))(domain, currentIP, entryID)
				} else {
					upLogger.Printf("No need to update DNS v4 record for %s", domain)
				}
			}
		}
		if Config.V6.Enabled {
			if !v6 {
				upLogger.Print("No IPv6 connectivity, skipping...")
				continue
			}
			//Check current IP
			currentIP := getCurrentIP(6)
			if currentIP == "" {
				upLogger.Printf("No IPv6 found, skipping...")
				continue
			}
			// Get DNS IP
			force, dnsIP, entryID := dns.GetV6.(func(string) (bool, string, string))(domain)
			//  Update DNS IP if needed
			if force {
				upLogger.Printf("Forcing update of DNS record for %s (%s)...", domain, entryID)
				dns.UpdateV6.(func(string, string, string))(domain, currentIP, entryID)
			} else {
				if dnsIP == "" {
					// Create DNS
					upLogger.Printf("No DNS record found for %s, creating...", domain)
					dns.CreateV6.(func(string, string))(domain, currentIP)
				} else if currentIP != dnsIP {
					// Update DNS
					upLogger.Printf("Updating DNS record for %s (%s)...", domain, entryID)
					dns.UpdateV6.(func(string, string, string))(domain, currentIP, entryID)
				} else {
					// No need to update DNS record for domain
					upLogger.Printf("No need to update DNS v6 record for %s", domain)
				}
			}
		}
	}
	upLogger.Print("Updater finished!")
}

func getCurrentIP(v int) string {
    var checkURL string
    var version string

    if v == 4 {
        checkURL = Config.V4.Check_url
        version = "v4"
    } else if v == 6 {
        checkURL = Config.V6.Check_url
        version = "v6"
    } else {
        upLogger.Printf("Invalid IP version requested: %d", v)
        return ""
    }

    //Check current IP
    cli := &http.Client{
        Timeout: 10 * time.Second, // Add timeout to prevent hanging
    }
    IPreq, err := http.NewRequest("GET", checkURL, nil)
    if err != nil {
        upLogger.Printf("Error creating request for %s IP check: %v", version, err)
        return ""
    }

    resp, err := cli.Do(IPreq)
    if err != nil {
        upLogger.Printf("Error getting %s IP: %v", version, err)
        return ""
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        upLogger.Printf("%s IP check returned non-OK status: %d", version, resp.StatusCode)
        return ""
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        upLogger.Printf("Error reading %s IP response: %v", version, err)
        return ""
    }

    ip := strings.TrimSpace(string(body))
    upLogger.Printf("Current %s IP: %s", version, ip)
    return ip
}
