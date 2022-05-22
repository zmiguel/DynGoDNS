package configurator

import (
	"log"
	"os"

	"github.com/zmiguel/DynGoDNS/internal/types"
	"gopkg.in/yaml.v3"
)

func ReadConfig(cfg *types.Config, filename string) {
	confLogger := log.New(os.Stdout, "[Config]     ", log.LstdFlags)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)

	if err != nil {
		confLogger.Fatal(err)
	}

	confLogger.Print("Configuration file loaded!")
}
