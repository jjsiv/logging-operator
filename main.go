package main

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type Input interface {
	Input()
}

type Tail struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func (t Tail) Input() {}

type MyStruct struct {
	Inputs []MyStructInput `yaml:"inputs"`
}

type MyStructInput struct {
	Name  string `yaml:"name"`
	Input `yaml:",inline"`
}

func main() {
	tail := Tail{
		Name: "tail",
		Path: "/tmp",
	}
	m := MyStruct{
		Inputs: []MyStructInput{
			{
				Name:  "tail",
				Input: tail,
			},
		},
	}

	b, err := yaml.Marshal(m)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(string(b))
}
