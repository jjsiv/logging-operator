package output

import (
	"github.com/jjsiv/logging-operator/internal/fluentbit"
	"github.com/jjsiv/logging-operator/internal/fluentbit/outputs"
)

type Stdout struct{}

func (s *Stdout) ToFluentBitOutput(opts *BuildOptions) fluentbit.Output {
	return outputs.Stdout{
		Tag: opts.Tag,
	}
}
