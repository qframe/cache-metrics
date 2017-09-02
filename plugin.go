package qcache_metrics

import (
	"fmt"
	"net/http"
	"time"
	"github.com/urfave/negroni"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/constants"
)

const (
	version   = "0.0.0"
	pluginTyp = qtypes_constants.CACHE
	pluginPkg = "metrics"
)


type Plugin struct {
	*qtypes_plugin.Plugin
}



func New(b qtypes_plugin.Base, name string) (Plugin, error) {
	p := qtypes_plugin.NewNamedPlugin(b, pluginTyp, pluginPkg, name, version)
	return Plugin{
		Plugin: p,
	}, nil
}



func (p *Plugin) Run() (err error){
	p.Log("notice", fmt.Sprintf("Start plugin v%s", p.Version))
	dc := p.QChan.Data.Join()
	done := p.QChan.Done.Join()
	go p.startHTTP()
	for {
		select {
		case <-dc.Read:
			continue
		case <- done.Read:
			break
		case err := <- p.ErrChan:
			return err
		}

	}
	return
}

func (p *Plugin) startHTTP() {
	mux := http.NewServeMux()
	n := negroni.New()
	n.UseHandler(mux)
	n.Use(negroni.HandlerFunc(p.LogMiddleware))
	addr := fmt.Sprintf("0.0.0.0:8124")
	p.Log("info", fmt.Sprintf("Start metrics-endpoint: %s", addr))
	p.ErrChan <- http.ListenAndServe(addr, n)
}

func (p *Plugin) LogMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	now := time.Now()
	next(rw, r)
	dur := time.Now().Sub(now)
	p.Log("trace", fmt.Sprintf("%s took %s", r.URL.String(), dur.String()))
}