package output

import "github.com/jjsiv/logging-operator/internal/fluentbit"

type Output interface {
	ToFluentBitOutput() *fluentbit.Output
}
