package handler

import (
	"net/http"

	"github.com/DazWilkin/Aetos/opts"

	"google.golang.org/protobuf/encoding/protojson"
)

type Varz struct {
	opts *opts.Aetos
}

func NewVarz(opts *opts.Aetos) *Varz {
	return &Varz{
		opts: opts,
	}
}

func (v *Varz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := protojson.Marshal(v.opts.Get())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
