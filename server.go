package main

import (
	"bitbucket.org/chrj/smtpd"
	"bytes"
	"fmt"
	"net/http"
)

type Server struct {
	*smtpd.Server
	hooks Hooks
}

func NewServer(srv *smtpd.Server, hooks Hooks) *Server {
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

	return smtpd.Error{
		Code:    550,
		Message: "address rejected: user unknown in local recipient table",
	}
}

func (s *Server) ServeSMTP(_ smtpd.Peer, env smtpd.Envelope) error {
	msg, err := NewMessage(bytes.NewBuffer(env.Data))

	if err != nil {
		fmt.Println("could not parse mail:", err)

		return smtpd.Error{
			Code:    554,
			Message: "mail rejected: invalid format",
		}
	}

	for _, addr := range env.Recipients {
		hook, ok := s.hooks.Find(addr)

		if !ok {
			fmt.Println("could not find hook for address:", addr)

			return smtpd.Error{
				Code:    451,
				Message: "internal server error",
			}
		}
		buf := bytes.NewBuffer(nil)
		enc := NewMultipartEncoder(buf)

		if err := enc.Encode("mail", NewMail(env.Sender, addr, msg)); err != nil {
			fmt.Println("could not encode request body:", err)

			return smtpd.Error{
				Code:    451,
				Message: "internal server error",
			}
		}
		enc.Close()

		resp, err := http.Post(hook.Hook, enc.FormDataContentType(), buf)

		if err != nil {
			fmt.Println("could not dispatch message:", err)

			return smtpd.Error{
				Code:    451,
				Message: "unable to reach destination system",
			}
		}

		if 400 <= resp.StatusCode && resp.StatusCode <= 499 {
			fmt.Println("could not dispatch message: server responded with:", resp.Status)

			return smtpd.Error{
				Code:    554,
				Message: "destination system does not accept message",
			}
		}

		if resp.StatusCode < 200 || 299 < resp.StatusCode {
			fmt.Println("could not dispatch message: server responded with:", resp.Status)

			return smtpd.Error{
				Code:    451,
				Message: "internal server error",
			}
		}
		fmt.Println("relayed mail for", addr, "to", hook.Hook)
	}
	return nil
}
