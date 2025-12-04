package outputs

type Stdout struct {
	Tag string `json:"tag" yaml:"tag"`
}

func (s Stdout) OutputName() string {
	return "stdout"
}
