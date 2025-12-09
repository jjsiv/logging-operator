package input

import (
	"fmt"
	"strings"

	"github.com/jjsiv/logging-operator/internal/fluentbit"
	"github.com/jjsiv/logging-operator/internal/fluentbit/inputs"
)

// Tail input plugin tails logs of pods in a namespace.
// +kubebuilder:object:generate=true
type Tail struct {
	// Tail logs only from the specified pods. Pods not included are implicitly excluded.
	IncludePods []Pod `json:"includePods,omitempty" yaml:"pod,omitempty"`
	// Exclude specified pods from tailing.
	ExcludePods   []Pod  `json:"excludePods,omitempty" yaml:"pod,omitempty"`
	ReadFromHead  bool   `json:"readFromHead,omitempty" yaml:"readFromHead,omitempty"`
	BufferMaxSize string `json:"bufferMaxSize,omitempty" yaml:"bufferMaxSize,omitempty"`
}

// TODO: we need to set tag and path
func (t *Tail) ToFluentBitInput(opts *BuildOptions) fluentbit.Input {
	logRoot := "/var/log/pods"
	logFiles := "*.log"

	var logPaths []string
	if len(t.IncludePods) == 0 {
		logPaths = append(logPaths, fmt.Sprintf("%s/%s_*/*/%s", logRoot, opts.Namespace, logFiles))
	}

	for _, pod := range t.IncludePods {
		if len(pod.Containers) == 0 {
			logPaths = append(logPaths, fmt.Sprintf("%s/%s_%s_*/*/%s", logRoot, opts.Namespace, pod.Name, logFiles))
		}
		for _, ctr := range pod.Containers {
			logPaths = append(logPaths, fmt.Sprintf("%s/%s_%s_*/%s/%s", logRoot, opts.Namespace, pod.Name, ctr, logFiles))
		}
	}

	var excludePaths []string
	for _, pod := range t.ExcludePods {
		if len(pod.Containers) == 0 {
			excludePaths = append(logPaths, fmt.Sprintf("%s/%s_%s_*/*/%s", logRoot, opts.Namespace, pod.Name, logFiles))
		}
		for _, ctr := range pod.Containers {
			excludePaths = append(logPaths, fmt.Sprintf("%s/%s_%s_*/%s/%s", logRoot, opts.Namespace, pod.Name, ctr, logFiles))
		}
	}

	tail := inputs.Tail{
		Tag:           opts.Tag,
		Path:          strings.Join(logPaths, ","),
		ExcludePath:   strings.Join(excludePaths, ","),
		ReadFromHead:  t.ReadFromHead,
		BufferMaxSize: t.BufferMaxSize,
	}

	return tail
}

// +kubebuilder:object:generate=true
type Pod struct {
	// Name of the pod to exclude or include logs from. Wildcards are supported.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Container names to exclude or include.
	Containers []string `json:"containers,omitempty" yaml:"containers,omitempty"`
}
