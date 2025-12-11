package inputs

type Tail struct {
	Tag           string `json:"tag" yaml:"tag"`
	Path          string `json:"path" yaml:"path"`
	ExcludePath   string `json:"exclude_path,omitempty" yaml:"exclude_path,omitempty"`
	ReadFromHead  bool   `json:"read_from_head,omitempty" yaml:"read_from_head,omitempty"`
	BufferMaxSize string `json:"buffer_max_size,omitempty" yaml:"buffer_max_size,omitempty"`
}

func (t Tail) InputName() string {
	return "tail"
}
