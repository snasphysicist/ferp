package forward

import (
	"github.com/go-chi/chi/v5"
	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/log"
	"github.com/snasphysicist/ferp/v2/pkg/proxy"
	"github.com/snasphysicist/ferp/v2/pkg/url"
)

// Configure sets up on the router all proxy routes defined in the incomings
func Configure(r *chi.Mux, incs []configuration.Incoming) {
	for _, i := range incs {
		rm := proxy.Proxy{
			BaseURL: url.BaseURL{
				Protocol: i.Downstream.Protocol,
				Host:     i.Downstream.Host,
				Port:     i.Downstream.Port,
				Path:     i.Downstream.Base,
			},
			Mapper: i.Downstream.Mapper.Map,
		}
		log.Infof("For Incoming %#v constructed Remapper %#v ", i, rm)
		for _, mr := range i.MethodRouters {
			log.Infof("Configuring remapper %#v for incoming '%s' with %#v",
				rm, i.Path, mr)
			mr.Route(r, i.Path, rm.ForwardRequest)
		}
	}
}
