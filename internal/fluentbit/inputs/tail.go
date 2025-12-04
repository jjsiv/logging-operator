package inputs

type Tail struct {
	Tag           string `json:"tag" yaml:"tag"`
	Path          string `json:"path" yaml:"path"`
	ReadFromHead  bool   `json:"read_from_head" yaml:"read_from_head"`
	BufferMaxSize string `json:"buffer_max_size" yaml:"buffer_max_size"`
}

func (t Tail) InputName() string {
	return "tail"
}
