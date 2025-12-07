package input

import "github.com/jjsiv/logging-operator/internal/fluentbit"

type Input interface {
	ToFluentBitInput(*BuildOptions) *fluentbit.Input
}

type BuildOptions struct {
	Tag       string
	Namespace string
}
