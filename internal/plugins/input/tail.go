package input

import (
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

func (t *Tail) ToFluentBitInput() fluentbit.Input {
	tail := inputs.Tail{
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
