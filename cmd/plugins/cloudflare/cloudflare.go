package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zmiguel/DynGoDNS/internal/types"
)

var (
	api_endpoint = "https://api.cloudflare.com/client/v4"
	Config       types.Config
	domains_data []domains
	cloudLogger  = log.New(os.Stdout, "[Cloudflare] ", log.LstdFlags)
)

func Info() string {
	return "Cloudflare v1.0.0"
}

func Initialise() {
	if domains_data == nil {
		getDomainsData()
	}
}

func getDomainsData() {
	for _, domain := range Config.Domains {
		var data domains
		// check if domain contains a comma
		if strings.Contains(domain, ",") {
			var domain_clean string = strings.Split(domain, ",")[0]
			if len(strings.Split(domain, ",")) > 1 && strings.Split(domain, ",")[1] == "proxy" {
				data.proxy = true
			} else {
				data.proxy = false
			}
			data.domain = domain_clean
		} else {
			data.domain = domain
			data.proxy = false
		}
		data.zone = strings.Split(domain, ".")[len(strings.Split(domain, "."))-2] + "." + strings.Split(domain, ".")[len(strings.Split(domain, "."))-1]

		cli := &http.Client{}
		req, err := http.NewRequest("GET", api_endpoint+"/zones?name="+data.zone, nil)
		if err != nil {
			cloudLogger.Fatal(err)
		}
		if Config.DNS.Username != "" {
			req.Header = map[string][]string{
				"X-Auth-Email": {Config.DNS.Username},
				"X-Auth-Key":   {Config.DNS.Password},
				"Content-Type": {"application/json"},
			}
		} else {
			req.Header = map[string][]string{
				"X-Auth-User-Service-Key": {Config.DNS.Password},
				"Content-Type":            {"application/json"},
			}
		}

		resp, err := cli.Do(req)
		if err != nil {
			cloudLogger.Fatal(err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			cloudLogger.Fatal(err)
		}

		var result zones
		json.Unmarshal(body, &result)

		if result.Success && len(result.Result) > 0 {
			data.zone_id = result.Result[0].ID
		} else {
			data.zone_id = ""
		}

		domains_data = append(domains_data, data)
		cloudLogger.Printf("Entry: %+v", data)
	}
}

func getDomain(dom string) domains {
	for _, domain := range domains_data {
		if domain.domain == dom {
			return domain
		}
	}
	return domains{}
}

func GetV4(dom string) (bool, string, string) {
	// get current DNS records
	domain := getDomain(dom)
	cli := &http.Client{}
	DNSreq, err := http.NewRequest("GET", api_endpoint+"/zones/"+domain.zone_id+"/dns_records?type=A&name="+dom, nil)
	if err != nil {
		cloudLogger.Fatal(err)
	}
	if Config.DNS.Username != "" {
		DNSreq.Header = map[string][]string{
			"X-Auth-Email": {Config.DNS.Username},
			"X-Auth-Key":   {Config.DNS.Password},
			"Content-Type": {"application/json"},
		}
	} else {
		DNSreq.Header = map[string][]string{
			"X-Auth-User-Service-Key": {Config.DNS.Password},
			"Content-Type":            {"application/json"},
		}
	}

	resp, err := cli.Do(DNSreq)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	var result listDNS
	json.Unmarshal(body, &result)

	DNSip := ""
	if result.Success && len(result.Result) > 0 && result.Result[0].Name == dom {
		DNSip = result.Result[0].Content
		cloudLogger.Printf("Current v4 DNS: %s (Proxy: %t)", DNSip, result.Result[0].Proxied)
		if result.Result[0].Proxied != domain.proxy {
			cloudLogger.Printf("Proxy status mismatch! Expected: %t, Got: %t", domain.proxy, result.Result[0].Proxied)
			return true, DNSip, result.Result[0].ID
		}
		return false, DNSip, result.Result[0].ID
	}
	return false, DNSip, ""
}

func CreateV4(dom string, ip string) {
	// Create new DNS record
	domain := getDomain(dom)
	cli := &http.Client{}
	Addreq, err := http.NewRequest("POST", api_endpoint+"/zones/"+domain.zone_id+"/dns_records", nil)
	if err != nil {
		cloudLogger.Fatal(err)
	}
	if Config.DNS.Username != "" {
		Addreq.Header = map[string][]string{
			"X-Auth-Email": {Config.DNS.Username},
			"X-Auth-Key":   {Config.DNS.Password},
			"Content-Type": {"application/json"},
		}
	} else {
		Addreq.Header = map[string][]string{
			"X-Auth-User-Service-Key": {Config.DNS.Password},
			"Content-Type":            {"application/json"},
		}
	}

	Addreq.Body = io.NopCloser(strings.NewReader(`{"type":"A","name":"` + domain.domain + `","content":"` + ip + `","ttl":1,"proxied":` + strconv.FormatBool(domain.proxy) + `}`))

	cloudLogger.Printf("Setting domain: %s to %s proxy: %t", domain.domain, ip, domain.proxy)

	resp, err := cli.Do(Addreq)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		cloudLogger.Fatal(err)
	}

	var result modifyDNS
	json.Unmarshal(body, &result)

	if result.Success && result.Result.Content == ip {
		cloudLogger.Print("DNS record created")
	}
}

func UpdateV4(dom string, ip string, id string) {
	// Update DNS record
	domain := getDomain(dom)
	cli := &http.Client{}
	Addreq, err := http.NewRequest("PATCH", api_endpoint+"/zones/"+domain.zone_id+"/dns_records/"+id, nil)
	if err != nil {
		cloudLogger.Fatal(err)
	}
	if Config.DNS.Username != "" {
		Addreq.Header = map[string][]string{
			"X-Auth-Email": {Config.DNS.Username},
			"X-Auth-Key":   {Config.DNS.Password},
			"Content-Type": {"application/json"},
		}
	} else {
		Addreq.Header = map[string][]string{
			"X-Auth-User-Service-Key": {Config.DNS.Password},
			"Content-Type":            {"application/json"},
		}
	}

	Addreq.Body = io.NopCloser(strings.NewReader(`{"type":"A","name":"` + domain.domain + `","content":"` + ip + `","ttl":1,"proxied":` + strconv.FormatBool(domain.proxy) + `}`))

	cloudLogger.Printf("Setting domain: %s to %s proxy: %t", domain.domain, ip, domain.proxy)

	resp, err := cli.Do(Addreq)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	var result modifyDNS
	json.Unmarshal(body, &result)

	if result.Success && result.Result.Content == ip {
		cloudLogger.Print("DNS record updated")
	}
}

func GetV6(dom string) (bool, string, string) {
	// get current DNS records
	domain := getDomain(dom)
	cli := &http.Client{}
	DNSreq, err := http.NewRequest("GET", api_endpoint+"/zones/"+domain.zone_id+"/dns_records?type=AAAA&name="+dom, nil)
	if err != nil {
		cloudLogger.Fatal(err)
	}
	if Config.DNS.Username != "" {
		DNSreq.Header = map[string][]string{
			"X-Auth-Email": {Config.DNS.Username},
			"X-Auth-Key":   {Config.DNS.Password},
			"Content-Type": {"application/json"},
		}
	} else {
		DNSreq.Header = map[string][]string{
			"X-Auth-User-Service-Key": {Config.DNS.Password},
			"Content-Type":            {"application/json"},
		}
	}

	resp, err := cli.Do(DNSreq)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	var result listDNS
	json.Unmarshal(body, &result)

	DNSip := ""
	if result.Success && len(result.Result) > 0 && result.Result[0].Name == dom {
		DNSip = result.Result[0].Content
		cloudLogger.Printf("Current v6 DNS: %s (Proxy: %t)", DNSip, result.Result[0].Proxied)
		if result.Result[0].Proxied != domain.proxy {
			cloudLogger.Printf("Proxy status mismatch! Expected: %t, Got: %t", domain.proxy, result.Result[0].Proxied)
			return true, DNSip, result.Result[0].ID
		}
		return false, DNSip, result.Result[0].ID
	}
	return false, DNSip, ""
}

func CreateV6(dom string, ip string) {
	// Create new DNS record
	domain := getDomain(dom)
	cli := &http.Client{}
	Addreq, err := http.NewRequest("POST", api_endpoint+"/zones/"+domain.zone_id+"/dns_records", nil)
	if err != nil {
		cloudLogger.Fatal(err)
	}
	if Config.DNS.Username != "" {
		Addreq.Header = map[string][]string{
			"X-Auth-Email": {Config.DNS.Username},
			"X-Auth-Key":   {Config.DNS.Password},
			"Content-Type": {"application/json"},
		}
	} else {
		Addreq.Header = map[string][]string{
			"X-Auth-User-Service-Key": {Config.DNS.Password},
			"Content-Type":            {"application/json"},
		}
	}

	Addreq.Body = io.NopCloser(strings.NewReader(`{"type":"AAAA","name":"` + domain.domain + `","content":"` + ip + `","ttl":1,"proxied":` + strconv.FormatBool(domain.proxy) + `}`))

	cloudLogger.Printf("Setting domain: %s to %s proxy: %t", domain.domain, ip, domain.proxy)

	resp, err := cli.Do(Addreq)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	var result modifyDNS
	json.Unmarshal(body, &result)

	if result.Success && result.Result.Content == ip {
		cloudLogger.Print("DNS record created")
	}
}

func UpdateV6(dom string, ip string, id string) {
	// Update DNS record
	domain := getDomain(dom)
	cli := &http.Client{}
	Addreq, err := http.NewRequest("PATCH", api_endpoint+"/zones/"+domain.zone_id+"/dns_records/"+id, nil)
	if err != nil {
		cloudLogger.Fatal(err)
	}
	if Config.DNS.Username != "" {
		Addreq.Header = map[string][]string{
			"X-Auth-Email": {Config.DNS.Username},
			"X-Auth-Key":   {Config.DNS.Password},
			"Content-Type": {"application/json"},
		}
	} else {
		Addreq.Header = map[string][]string{
			"X-Auth-User-Service-Key": {Config.DNS.Password},
			"Content-Type":            {"application/json"},
		}
	}

	Addreq.Body = io.NopCloser(strings.NewReader(`{"type":"AAAA","name":"` + domain.domain + `","content":"` + ip + `","ttl":1,"proxied":` + strconv.FormatBool(domain.proxy) + `}`))

	cloudLogger.Printf("Setting domain: %s to %s proxy: %t", domain.domain, ip, domain.proxy)

	resp, err := cli.Do(Addreq)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cloudLogger.Fatal(err)
	}

	var result modifyDNS
	json.Unmarshal(body, &result)

	if result.Success && result.Result.Content == ip {
		cloudLogger.Print("DNS record updated")
	}
}

type domains struct {
	domain  string
	zone    string
	zone_id string
	proxy   bool
}

type zones struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []struct {
		ID                  string    `json:"id"`
		Name                string    `json:"name"`
		DevelopmentMode     int       `json:"development_mode"`
		OriginalNameServers []string  `json:"original_name_servers"`
		OriginalRegistrar   string    `json:"original_registrar"`
		OriginalDnshost     string    `json:"original_dnshost"`
		CreatedOn           time.Time `json:"created_on"`
		ModifiedOn          time.Time `json:"modified_on"`
		ActivatedOn         time.Time `json:"activated_on"`
		Owner               struct {
			ID struct {
			} `json:"id"`
			Email struct {
			} `json:"email"`
			Type string `json:"type"`
		} `json:"owner"`
		Account struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Permissions []string `json:"permissions"`
		Plan        struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Price        int    `json:"price"`
			Currency     string `json:"currency"`
			Frequency    string `json:"frequency"`
			LegacyID     string `json:"legacy_id"`
			IsSubscribed bool   `json:"is_subscribed"`
			CanSubscribe bool   `json:"can_subscribe"`
		} `json:"plan"`
		PlanPending struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Price        int    `json:"price"`
			Currency     string `json:"currency"`
			Frequency    string `json:"frequency"`
			LegacyID     string `json:"legacy_id"`
			IsSubscribed bool   `json:"is_subscribed"`
			CanSubscribe bool   `json:"can_subscribe"`
		} `json:"plan_pending"`
		Status      string   `json:"status"`
		Paused      bool     `json:"paused"`
		Type        string   `json:"type"`
		NameServers []string `json:"name_servers"`
	} `json:"result"`
}

type listDNS struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []struct {
		ID         string    `json:"id"`
		Type       string    `json:"type"`
		Name       string    `json:"name"`
		Content    string    `json:"content"`
		Proxiable  bool      `json:"proxiable"`
		Proxied    bool      `json:"proxied"`
		TTL        int       `json:"ttl"`
		Locked     bool      `json:"locked"`
		ZoneID     string    `json:"zone_id"`
		ZoneName   string    `json:"zone_name"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
		Data       struct {
		} `json:"data"`
		Meta struct {
			AutoAdded bool   `json:"auto_added"`
			Source    string `json:"source"`
		} `json:"meta"`
	} `json:"result"`
}

type modifyDNS struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   struct {
		ID         string    `json:"id"`
		Type       string    `json:"type"`
		Name       string    `json:"name"`
		Content    string    `json:"content"`
		Proxiable  bool      `json:"proxiable"`
		Proxied    bool      `json:"proxied"`
		TTL        int       `json:"ttl"`
		Locked     bool      `json:"locked"`
		ZoneID     string    `json:"zone_id"`
		ZoneName   string    `json:"zone_name"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
		Data       struct {
		} `json:"data"`
		Meta struct {
			AutoAdded bool   `json:"auto_added"`
			Source    string `json:"source"`
		} `json:"meta"`
	} `json:"result"`
}
