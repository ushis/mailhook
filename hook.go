package main

import (
	"fmt"
	"net/mail"
	"strings"
)

type Hooks []*Hook

func (hs Hooks) Find(addr string) (*Hook, bool) {
	addr, err := normalizeAddress(addr)

	if err != nil {
		return nil, false
	}

	for _, hook := range hs {
		if hook.includesAddr(addr) {
			return hook, true
		}
	}
	return nil, false
}

type Hook struct {
	Emails []string `yaml:"emails"`
	Hook   string   `yaml:"hook"`
}

func (h *Hook) includesAddr(addr string) bool {
	for _, email := range h.Emails {
		if addr == email {
			return true
		}

		if strings.HasPrefix(email, "@") && strings.HasSuffix(addr, email) {
			return true
		}

		if email == "@" {
			return true
		}
	}
	return false
}

func normalizeAddress(addr string) (string, error) {
	email, err := mail.ParseAddress(addr)

	if err != nil {
		return "", err
	}
	local, domain, err := splitAddress(email.Address)

	if err != nil {
		return "", err
	}

	if parts := strings.SplitN(local, "+", 2); len(parts) > 1 {
		local = parts[0]
	}
	return fmt.Sprintf("%s@%s", local, domain), nil
}

func splitAddress(addr string) (string, string, error) {
	parts := strings.SplitN(addr, "@", 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid address: %s", addr)
	}
	return parts[0], parts[1], nil
}
