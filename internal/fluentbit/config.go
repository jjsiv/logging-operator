package fluentbit

type Input interface {
	InputName() string
}

type Filter interface {
	FilterName() string
}

type Output interface {
	OutputName() string
}

type FluentBitConfig struct {
	Includes []string                     `json:"includes,omitempty" yaml:"includes,omitempty"`
	Pipeline *FluentBitConfigPipelineSpec `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
	Service  *FluentBitConfigServiceSpec  `json:"service,omitempty" yaml:"service,omitempty"`
}

type FluentBitConfigPipelineSpec struct {
	Inputs  []FluentBitPipelineInputSpec  `json:"inputs" yaml:"inputs"`
	Filters []FluentBitPipelineFilterSpec `json:"filters" yaml:"filters"`
	Outputs []FluentBitPipelineOutputSpec `json:"outputs" yaml:"outputs"`
}

type FluentBitConfigServiceSpec struct {
	Flush    float64 `json:"flush,omitempty" yaml:"flush,omitempty"`
	Grace    int64   `json:"grace,omitempty" yaml:"grace,omitempty"`
	LogLevel string  `json:"log_level,omitempty" yaml:"log_level,omitempty"`
}

type FluentBitPipelineInputSpec struct {
	Name  string `json:"name" yaml:"name"`
	Input `json:",inline" yaml:",inline"`
}

type FluentBitPipelineFilterSpec struct {
	Name   string `json:"name" yaml:"name"`
	Filter `json:",inline" yaml:",inline"`
}

type FluentBitPipelineOutputSpec struct {
	Name   string `json:"name" yaml:"name"`
	Output `json:",inline" yaml:",inline"`
}
