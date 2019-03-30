package heldiamgo

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

type HandleFunc func()

func PprofHttpHandler() (string, http.Handler) {
	return "/debug/pprof/*", handler("/debug/pprof/")
}

type handler string

func (name handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, string(name))
	path = strings.TrimLeft(path, "/")
	switch path {
	default:
		pprof.Index(w, r)
	case "cmdline":
		pprof.Cmdline(w, r)
	case "profile":
		pprof.Profile(w, r)
	case "symbol":
		pprof.Symbol(w, r)
	case "trace":
		pprof.Trace(w, r)
	}
}
