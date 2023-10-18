package main

import (
	"flag"
	"log"
	"os"
	"plugin"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/zmiguel/DynGoDNS/internal/configurator"
	"github.com/zmiguel/DynGoDNS/internal/types"
	"github.com/zmiguel/DynGoDNS/internal/updater"
)

var (
	cfg_file   string
	dns        types.DNS
	mainLogger = log.New(os.Stdout, "[DynGoDNS]   ", log.LstdFlags)
)

func init() {
	mainLogger.SetPrefix("[DynGoDNS]   ")
	flag.StringVar(&cfg_file, "config", "config.yaml", "Configuration file")
}

func main() {
	flag.Parse()

	var config types.Config

	mainLogger.Print("DynGoDNS v1.0.0")

	configurator.ReadConfig(&config, cfg_file)
	mainLogger.Printf("Detected DNS provider: %s", config.DNS.Provider)
	mainLogger.Print("Attempting to load plugin...")
	wd, wderr := os.Getwd()
	if wderr != nil {
		mainLogger.Fatal(wderr)
	}
	p, err := plugin.Open(wd + "/plugins/" + config.DNS.Provider + ".so")
	if err != nil {
		mainLogger.Fatal(err)
	}

	loadPluginFunctions(*p)
	updater.Initialise(dns, config)

	mainLogger.Printf("Loaded %s", dns.Info.(func() string)())
	mainLogger.Print("Configuring Plugin...")
	*dns.Config.(*types.Config) = config

	mainLogger.Print("Configuring Scheduler...")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(config.Check_interval).Do(updater.Update)

	mainLogger.Print("Starting Scheduler!")
	scheduler.StartBlocking()
}

func loadPluginFunctions(plug plugin.Plugin) {
	var err error
	// Info
	dns.Info, err = plug.Lookup("Info")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Config
	dns.Config, err = plug.Lookup("Config")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Initialise
	dns.Initialise, err = plug.Lookup("Initialise")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Get V4
	dns.GetV4, err = plug.Lookup("GetV4")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Get V6
	dns.GetV6, err = plug.Lookup("GetV6")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Update V4
	dns.UpdateV4, err = plug.Lookup("UpdateV4")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Update V6
	dns.UpdateV6, err = plug.Lookup("UpdateV6")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Create V4
	dns.CreateV4, err = plug.Lookup("CreateV4")
	if err != nil {
		mainLogger.Fatal(err)
	}
	// Create V6
	dns.CreateV6, err = plug.Lookup("CreateV6")
	if err != nil {
		mainLogger.Fatal(err)
	}
}
