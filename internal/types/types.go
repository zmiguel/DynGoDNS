package types

import "plugin"

type Config struct {
	Check_interval string   `yaml:"check_interval"`
	Domains        []string `yaml:"domains"`
	DNS            struct {
		Provider string `yaml:"provider"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		OPT1     string `yaml:"opt1"`
		OPT2     string `yaml:"opt2"`
	} `yaml:"dns"`
	V4 struct {
		Enabled   bool   `yaml:"enabled"`
		Delete	  bool   `yaml:"delete"`
		Check_url string `yaml:"check_url"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"v4"`
	V6 struct {
		Enabled   bool   `yaml:"enabled"`
		Delete	  bool   `yaml:"delete"`
		Check_url string `yaml:"check_url"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"v6"`
}

type DNS struct {
	Info       plugin.Symbol
	Config     plugin.Symbol
	Initialise plugin.Symbol
	GetV4      plugin.Symbol
	GetV6      plugin.Symbol
	UpdateV4   plugin.Symbol
	UpdateV6   plugin.Symbol
	CreateV4   plugin.Symbol
	CreateV6   plugin.Symbol
	DeleteV4   plugin.Symbol
	DeleteV6   plugin.Symbol
}
