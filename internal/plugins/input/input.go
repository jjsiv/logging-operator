package input

import "github.com/jjsiv/logging-operator/internal/fluentbit"

type Input interface {
	ToFluentBitInput() *fluentbit.Input
}
