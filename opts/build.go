package opts

type Build struct {
	Namespace string
	Subsystem string

	GitCommit string
	GoVersion string
	OSVersion string
	StartTime int64
}
