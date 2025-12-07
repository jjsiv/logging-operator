package input

import (
	"fmt"

	"github.com/jjsiv/logging-operator/internal/fluentbit"
	"github.com/jjsiv/logging-operator/internal/fluentbit/inputs"
)

// Tail input plugin tails logs of pods in a namespace.
type Tail struct {
	// Tail logs only from the specified pods. If empty, logs from all pods in the namespace will be collected.
	Pods          []Pod  `json:"pod,omitempty" yaml:"pod,omitempty"`
	ReadFromHead  bool   `json:"readFromHead,omitempty" yaml:"readFromHead,omitempty"`
	BufferMaxSize string `json:"bufferMaxSize,omitempty" yaml:"bufferMaxSize,omitempty"`
}

// TODO: we need to set tag and path
func (t *Tail) ToFluentBitInput(opts *BuildOptions) fluentbit.Input {
	logRoot := "/var/log/pods"
	logFiles := "*.log"

	var logPaths []string
	if len(t.Pods) == 0 {
		logPaths = append(logPaths, fmt.Sprintf("%s/%s_*/*/%s", logRoot, opts.Namespace, logFiles))
	} else {
		for _, pod := range t.Pods {

		}
	}
	tail := inputs.Tail{
		Path:          logRoot + "/*/*/*.log",
		ReadFromHead:  t.ReadFromHead,
		BufferMaxSize: t.BufferMaxSize,
	}

	return tail
}

type Pod struct {
	// Name of the pod to tail logs from.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// A list of containers to explicitly include. If unset, all containers will be used.
	Containers []string `json:"containers,omitempty" yaml:"containers,omitempty"`
}
