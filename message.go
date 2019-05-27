package main

import (
	"github.com/jhillyerd/enmime"
	"io"
	"net/mail"
	"strings"
	"time"
)

type Message struct {
	From        []*Address    `multipart:"from"`
	To          []*Address    `multipart:"to"`
	Cc          []*Address    `multipart:"cc"`
	Bcc         []*Address    `multipart:"bcc"`
	ReplyTo     []*Address    `multipart:"reply_to"`
	Subject     string        `multipart:"subject"`
	Date        *Time         `multipart:"date"`
	MessageID   string        `multipart:"message_id"`
	InReplyTo   string        `multipart:"in_reply_to"`
	References  []string      `multipart:"references"`
	Text        string        `multipart:"text"`
	HTML        string        `multipart:"html"`
	Attachments []*Attachment `multipart:"attachments"`
}

type Address struct {
	Name  string `multipart:"name"`
	Email string `multipart:"email"`
}

func NewMessage(r io.Reader) (*Message, error) {
	env, err := enmime.ReadEnvelope(r)

	if err != nil {
		return nil, err
	}

	msg := &Message{
		Subject:     env.GetHeader("Subject"),
		MessageID:   env.GetHeader("Message-ID"),
		InReplyTo:   env.GetHeader("In-Reply-To"),
		References:  strings.Fields(env.GetHeader("References")),
		Text:        env.Text,
		HTML:        env.HTML,
		Attachments: make([]*Attachment, len(env.Attachments)),
	}

	for i, attachment := range env.Attachments {
		msg.Attachments[i] = &Attachment{attachment}
	}

	if msg.From, err = readAddressListHeader(env, "From"); err != nil {
		return nil, err
	}

	if msg.To, err = readAddressListHeader(env, "To"); err != nil {
		return nil, err
	}

	if msg.Cc, err = readAddressListHeader(env, "Cc"); err != nil {
		return nil, err
	}

	if msg.Bcc, err = readAddressListHeader(env, "Bcc"); err != nil {
		return nil, err
	}

	if msg.ReplyTo, err = readAddressListHeader(env, "Reply-To"); err != nil {
		return nil, err
	}

	if msg.Date, err = readDateHeader(env); err != nil {
		return nil, err
	}
	return msg, nil
}

func readAddressListHeader(env *enmime.Envelope, key string) ([]*Address, error) {
	list, err := env.AddressList(key)

	if err != nil && err != mail.ErrHeaderNotPresent {
		return nil, err
	}
	emails := make([]*Address, len(list))

	for i, addr := range list {
		emails[i] = &Address{Name: addr.Name, Email: addr.Address}
	}
	return emails, nil
}

func readDateHeader(env *enmime.Envelope) (*Time, error) {
	hdr := env.GetHeader("Date")

	if hdr == "" {
		return &Time{time.Now()}, nil
	}
	date, err := mail.ParseDate(hdr)

	return &Time{date}, err
}
