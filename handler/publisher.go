package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DazWilkin/Aetos/api/v1alpha1"
	"github.com/DazWilkin/Aetos/opts"

	"google.golang.org/protobuf/encoding/protojson"
)

type Publisher struct {
	Config *opts.Aetos
}

func NewPublisher(config *opts.Aetos) *Publisher {
	return &Publisher{
		Config: config,
	}
}
func (p *Publisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		// Not ideal but there's no logger to use to handle the error
		_, _ = w.Write([]byte("expected POST"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Not ideal but there's no logger to use to handle the error
		_, _ = w.Write([]byte("unable to read request body"))
		return
	}

	// Request should unmarshal to pb.AetosPublishRequest
	rqst := &v1alpha1.AetosPublishRequest{}
	if err := protojson.Unmarshal(body, rqst); err != nil {
		// Expected request body to unmarshal
		w.WriteHeader(http.StatusBadRequest)
		// Not ideal but there's no logger to use to handle the error
		_, _ = w.Write([]byte("expected to be able to unmarshal request body"))
		return
	}

	// Validate AetosPublishRequest against Foo
	if len(rqst.Labels) > int(p.Config.NumLabels) {
		w.WriteHeader(http.StatusBadRequest)
		// Not ideal but there's no logger to use to handle the error
		_, _ = fmt.Fprintf(w, "Too many labels (max: %d)", p.Config.NumLabels)
		return
	}
	if len(rqst.Metrics) > int(p.Config.NumMetrics) {
		w.WriteHeader(http.StatusBadRequest)
		// Not ideal but there's no logger to use to handle the error
		_, _ = fmt.Fprintf(w, "Too many metrics (max: %d)", p.Config.NumMetrics)
		return
	}

	// Update Foo
	if err := p.Config.Set(rqst); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
