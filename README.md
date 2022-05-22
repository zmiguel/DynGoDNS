# DynGoDNS

Dynamic DNS updater written in GO

## Features

* Plugin based system for DNS providers
* V4 and V6 support
* Automatic update of DNS records in a timer

## Example configuration file

```yaml
check_interval: 1m # s m h
domains:
    - anything.here.com
    - another.domain.com
    - third.domain.com
dns:
    provider: cloudflare # must match plugin name
    username: 
    password: 
    opt1: null
    opt2: null
v4:
    enabled: true
    check_url: http://ip1.dynupdate.no-ip.com/
v6:
    enabled: true
    check_url: http://ip1.dynupdate6.no-ip.com/
```

## Installation

Download the binary for your platform and the pre-compiled plugin for your DNS provider

Place the plugin in a plugins folder in the same directory as your binary

Place your configuration file (`config.yaml`) in the same directory as the binary or use the `-config` flag to specify its location

## Contributing

There's a plugin template in the `plugin` folder.
the plugin **MUST** be named after the provider it's for and must return the same values as expected. the configuration file has 2 optional parameter in case your provider needs them.

Please check the included Cloudflare plugin for a more in depth example.
