package main

import (
	"path/filepath"
)

type HookDir struct {
	path string
}

func NewHookDir(path string) *HookDir {
	return &HookDir{path: path}
}

func (hd *HookDir) Hooks() (Hooks, error) {
	paths, err := filepath.Glob(filepath.Join(hd.path, "*.yml"))

	if err != nil {
		return nil, err
	}
	hooks := Hooks{}

	for _, path := range paths {
		fileHooks, err := NewHookFile(path).Hooks()

		if err != nil {
			return nil, err
		}
		hooks = append(hooks, fileHooks...)
	}
	return hooks, nil
}
