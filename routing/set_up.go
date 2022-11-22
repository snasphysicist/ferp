package routing

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"snas.pw/ferp/v2/configuration"
	"snas.pw/ferp/v2/proxy"
)

func SetUpRedirects(r *chi.Mux, rds []configuration.Redirect) {
	for _, rd := range rds {
		for _, mr := range rd.MethodRouters {
			mr(r, rd.From, redirector(rd.To))
			log.Printf("Configuring redirect from '%s' to '%s' with %#v",
				rd.From, rd.To, mr)
		}
	}
}

func redirector(to string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("location", to)
		w.WriteHeader(http.StatusFound)
	}
}

func SetUpForwarding(r *chi.Mux, incs []configuration.Incoming) {
	for _, i := range incs {
		rm := proxy.Remapper{
			Protocol: i.Downstream.Protocol,
			Host:     i.Downstream.Host,
			Port:     i.Downstream.Port,
			Base:     i.Downstream.Base,
			Mapper:   i.Downstream.Mapper.Map,
		}
		for _, mr := range i.MethodRouters {
			log.Printf("Configuring remapper %#v for incoming '%s' with %#v",
				rm, i.Path, mr)
			mr(r, i.Path, rm.ForwardRequest)
		}
	}
}
