module eos-arb

go 1.16

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.8 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/go-redis/redis/v8 v8.10.0
	github.com/gocarina/gocsv v0.0.0-20210516172204-ca9e8a8ddea8
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/panyanyany/eos-go v0.9.4-0.20210621092235-c9610ae7ce20 // indirect
	github.com/parnurzeal/gorequest v0.2.16
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.0
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
	moul.io/http2curl v1.0.0 // indirect
)

//replace github.com/panyanyany/eos-go v0.9.1 => /var/www/go/src/github.com/panyanyany/eos-go
