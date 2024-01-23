package handler

import (
	"fmt"
	"io"
	"net/http"

	pb "github.com/DazWilkin/Aetos/protos"
	"github.com/DazWilkin/Aetos/xxx"

	"google.golang.org/protobuf/encoding/protojson"
)

type Publisher struct {
	Config *xxx.Config
}

func NewPublisher(config *xxx.Config) *Publisher {
	return &Publisher{
		Config: config,
	}
}
func (p *Publisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected POST"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to read request body"))
		return
	}

	// Request should unmarshal to pb.AetosPublishRequest
	rqst := &pb.AetosPublishRequest{}
	if err := protojson.Unmarshal(body, rqst); err != nil {
		// Expected request body to unmarshal
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected to be able to unmarshal request body"))
		return
	}

	// Validate AetosPublishRequest against Foo
	if len(rqst.Labels) > int(p.Config.NumLabels) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Too many labels (max: %d)", p.Config.NumLabels)
		return
	}
	if len(rqst.Metrics) > int(p.Config.NumMetrics) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Too many metrics (max: %d)", p.Config.NumMetrics)
		return
	}

	// Update Foo
	p.Config.Set(rqst)

}
