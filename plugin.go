package qcache_metrics

import (
	"fmt"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/urfave/negroni"
	"net/http"
	"time"
)

const (
	version   = "0.0.0"
	pluginTyp = qtypes.CACHE
	pluginPkg = "metrics"
)


type Plugin struct {
	qtypes.Plugin
}



func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	p := qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version)
	return Plugin{
		Plugin: 			p,
	}, nil
}



func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start plugin v%s", p.Version))
	dc := p.QChan.Data.Join()
	go p.startHTTP()
	for {
		select {
		case <-dc.Read:
			continue
		}
	}
}

func (p *Plugin) startHTTP() {
	mux := http.NewServeMux()
	n := negroni.New()
	n.UseHandler(mux)
	n.Use(negroni.HandlerFunc(p.LogMiddleware))
	addr := fmt.Sprintf("0.0.0.0:8124")
	p.Log("info", fmt.Sprintf("Start metrics-endpoint: %s", addr))
	http.ListenAndServe(addr, n)
}

func (p *Plugin) LogMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	now := time.Now()
	next(rw, r)
	dur := time.Now().Sub(now)
	p.Log("trace", fmt.Sprintf("%s took %s", r.URL.String(), dur.String()))
}