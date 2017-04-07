package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type HookFile struct {
	path string
}

func NewHookFile(path string) *HookFile {
	return &HookFile{path: path}
}

func (hf *HookFile) Hooks() (Hooks, error) {
	f, err := os.Open(hf.path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}
	hooks := Hooks{}

	if err := yaml.Unmarshal(buf, &hooks); err != nil {
		return nil, err
	}
	return hooks, nil
}
