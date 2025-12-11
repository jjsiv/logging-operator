package output

import "github.com/jjsiv/logging-operator/internal/fluentbit"

type Output interface {
	ToFluentBitOutput(*BuildOptions) *fluentbit.Output
}

type BuildOptions struct {
	Tag       string
	Namespace string
}
