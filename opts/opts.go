package opts

type Opts struct {
	namespace string
	subsystem string
}

func NewOpts(namespace, subsystem string) *Opts {
	return &Opts{
		namespace: namespace,
		subsystem: subsystem,
	}
}

// NewBuildOpts is a method that creates new NewBuildOpts opts
func (o *Opts) NewBuildOpts(GitCommit, GoVersion, OSVersion string, StartTime int64) *Build {
	return &Build{
		Namespace: o.namespace,
		Subsystem: o.subsystem,
		GitCommit: GitCommit,
		GoVersion: GoVersion,
		OSVersion: OSVersion,
		StartTime: StartTime,
	}
}

// NewAetosOpts is a method that creates new NewAetosOpts opts
func (o *Opts) NewAetosOpts(cardinality, numLabels, numMetrics uint8) *Aetos {
	return &Aetos{
		namespace:   o.namespace,
		subsystem:   o.subsystem,
		cardinality: cardinality,
		NumLabels:   numLabels,
		NumMetrics:  numMetrics,
		labels:      []string{},
		metrics:     []string{},
	}
}
