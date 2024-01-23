package handler

import (
	"net/http"

	"github.com/DazWilkin/Aetos/xxx"

	"google.golang.org/protobuf/encoding/protojson"
)

type Varz struct {
	Config *xxx.Config
}

func NewVarz(config *xxx.Config) *Varz {
	return &Varz{
		Config: config,
	}
}

func (v *Varz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := protojson.Marshal(v.Config.Get())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
