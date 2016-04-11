package main

import (
	"bitbucket.org/chrj/smtpd"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	*smtpd.Server
	hooks *Hooks
}

func NewServer(srv *smtpd.Server, hooks *Hooks) *Server {
	s := &Server{srv, hooks}
	s.Server.Handler = s.ServeSMTP
	s.Server.RecipientChecker = s.CheckRecipient
	return s
}

func (s *Server) CheckRecipient(_ smtpd.Peer, addr string) error {
	if _, ok := s.hooks.Find(addr); ok {
		return nil
	}
	fmt.Println("denied unknown recipient address:", addr)
	return smtpd.Error{550, "address rejected: user unknown in local recipient table"}
}

func (s *Server) ServeSMTP(_ smtpd.Peer, env smtpd.Envelope) error {
	msg, err := NewMessage(bytes.NewBuffer(env.Data))

	if err != nil {
		fmt.Println("could not parse mail:", err)
		return smtpd.Error{565, "mail rejected: invalid format"}
	}

	for _, addr := range env.Recipients {
		hook, ok := s.hooks.Find(addr)

		if !ok {
			fmt.Println("could not find hook for address:", addr)
			return smtpd.Error{430, "internal server error"}
		}
		buf := bytes.NewBuffer(nil)

		if err := json.NewEncoder(buf).Encode(NewPayload(env.Sender, addr, msg)); err != nil {
			fmt.Println("could not json encode message:", err)
			return smtpd.Error{430, "internal server error"}
		}
		resp, err := http.Post(hook.Hook, "application/json", buf)

		if err != nil {
			fmt.Println("could not dispatch message:", err)
			return smtpd.Error{420, "internal server error"}
		}

		if resp.StatusCode < 200 || 299 < resp.StatusCode {
			fmt.Println("could not dispatch message: server responded with:", resp.Status)
			return smtpd.Error{420, "internal server error"}
		}
		fmt.Println("relayed mail for", addr, "to", hook.Hook)
	}
	return nil
}
