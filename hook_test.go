package main

import (
	"gopkg.in/yaml.v2"
	"testing"
)

const hooksYaml = `
---
- hook: https://example.com/v1/mails
  emails:
    - 'alice@example.com'
    - '@example.net'

- hook: https://example.com/v2/mails
  emails:
    - '@example.org'
    - 'bob@example.io'

- hook: https://example.com/v3/mails
  emails:
    - '@'
`

func parseHooks() (Hooks, error) {
	hooks := Hooks{}

	if err := yaml.Unmarshal([]byte(hooksYaml), &hooks); err != nil {
		return nil, err
	}
	return hooks, nil
}

func TestFindByAddress(t *testing.T) {
	hooks, err := parseHooks()

	if err != nil {
		t.Fatal(err)
	}
	hook, ok := hooks.Find("alice@example.com")

	if !ok {
		t.Fatal("expected to find hook")
	}

	if hook.Hook != "https://example.com/v1/mails" {
		t.Fatalf("\"https://example.com/v1/mails\" != \"%s\"", hook.Hook)
	}
	hook, ok = hooks.Find("bob@example.io")

	if !ok {
		t.Fatal("expected to find hook")
	}

	if hook.Hook != "https://example.com/v2/mails" {
		t.Fatalf("\"https://example.com/v2/mails\" != \"%s\"", hook.Hook)
	}
}

func TestFindByTaggedAddress(t *testing.T) {
	hooks, err := parseHooks()

	if err != nil {
		t.Fatal(err)
	}
	hook, ok := hooks.Find("alice+249c01d587797d3451b03f90a2ecbc7a@example.com")

	if !ok {
		t.Fatal("expected to find hook")
	}

	if hook.Hook != "https://example.com/v1/mails" {
		t.Fatalf("\"https://example.com/v1/mails\" != \"%s\"", hook.Hook)
	}
}

func TestFindByDomain(t *testing.T) {
	hooks, err := parseHooks()

	if err != nil {
		t.Fatal(err)
	}
	hook, ok := hooks.Find("42b71135b144c5c5a1c803c9dedd6ceb@example.net")

	if !ok {
		t.Fatal("expected to find hook")
	}

	if hook.Hook != "https://example.com/v1/mails" {
		t.Fatalf("\"https://example.com/v1/mails\" != \"%s\"", hook.Hook)
	}
	hook, ok = hooks.Find("a73bd5e153e6f1321ebaf62ad01a93c8@example.org")

	if !ok {
		t.Fatal("expected to find hook")
	}

	if hook.Hook != "https://example.com/v2/mails" {
		t.Fatalf("\"https://example.com/v1/mails\" != \"%s\"", hook.Hook)
	}
}

func TestFindByCatchAll(t *testing.T) {
	hooks, err := parseHooks()

	if err != nil {
		t.Fatal(err)
	}
	hook, ok := hooks.Find("32cf34012668@41248827644a.e1e86c9fa707")

	if !ok {
		t.Fatal("expected to find hook")
	}

	if hook.Hook != "https://example.com/v3/mails" {
		t.Fatalf("\"https://example.com/v3/mails\" != \"%s\"", hook.Hook)
	}
}

func TestFindNothing(t *testing.T) {
	hooks, err := parseHooks()

	if err != nil {
		t.Fatal(err)
	}
	hooks = hooks[:len(hooks)-1]

	if _, ok := hooks.Find("32cf34012668@41248827644a.e1e86c9fa707"); ok {
		t.Fatal("expected not to find anything")
	}
}
