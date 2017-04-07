package main

import (
	"path/filepath"
	"testing"
)

func TestHookDirHooks(t *testing.T) {
	hooks, err := NewHookDir(filepath.Join("testdata", "hooks")).Hooks()

	if err != nil {
		t.Fatal(err)
	}

	if len(hooks) != 3 {
		t.Fatalf("3 != %d", len(hooks))
	}
	hook, ok := hooks.Find("a@a.com")

	if !ok {
		t.Fatalf("couldn't find hook for a@a.com")
	}

	if hook.Hook != "http://a.com/a" {
		t.Fatalf("http://a.com/a != %s", hook.Hook)
	}
	hook, ok = hooks.Find("c@a.com")

	if !ok {
		t.Fatalf("couldn't find hook for c@a.com")
	}

	if hook.Hook != "http://a.com/b" {
		t.Fatalf("http://a.com/b != %s", hook.Hook)
	}
	hook, ok = hooks.Find("something@b.net")

	if !ok {
		t.Fatalf("couldn't find hook for something@b.net")
	}

	if hook.Hook != "http://b.net/" {
		t.Fatalf("http://b.net/ != %s", hook.Hook)
	}
}
