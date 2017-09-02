package main

import (
	"log"
	"github.com/zpatrick/go-config"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/cache-metrics"
	"github.com/qframe/types/plugin"
)

func main() {
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"log.level": "debug",
		//"log.only-plugins": "",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	b := qtypes_plugin.NewBase(qChan, cfg)
	p, err := qcache_metrics.New(b, "metrics")
	if err != nil {
		log.Fatalf("[EE] Failed to create cache: %v", err)
	}
	go p.Run()
	bg := p.QChan.Data.Join()
	for {
		select {
		case <- bg.Read:
			continue

		}
	}
}
