package main

import (
	"github.com/jhillyerd/go.enmime"
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
	msg, err := mail.ReadMessage(r)

	if err != nil {
		return nil, err
	}
	b, err := enmime.ParseMIMEBody(msg)

	if err != nil {
		return nil, err
	}

	m := &Message{
		Subject:     b.GetHeader("Subject"),
		MessageID:   b.GetHeader("Message-ID"),
		InReplyTo:   b.GetHeader("In-Reply-To"),
		References:  strings.Fields(b.GetHeader("References")),
		Text:        b.Text,
		HTML:        b.HTML,
		Attachments: make([]*Attachment, len(b.Attachments)),
	}

	for i, attachment := range b.Attachments {
		m.Attachments[i] = &Attachment{attachment}
	}

	if m.From, err = readAddressListHeader(b, "From"); err != nil {
		return nil, err
	}

	if m.To, err = readAddressListHeader(b, "To"); err != nil {
		return nil, err
	}

	if m.Cc, err = readAddressListHeader(b, "Cc"); err != nil {
		return nil, err
	}

	if m.Bcc, err = readAddressListHeader(b, "Bcc"); err != nil {
		return nil, err
	}

	if m.ReplyTo, err = readAddressListHeader(b, "Reply-To"); err != nil {
		return nil, err
	}

	if m.Date, err = readDateHeader(msg.Header); err != nil {
		return nil, err
	}
	return m, nil
}

func readAddressListHeader(m *enmime.MIMEBody, key string) ([]*Address, error) {
	list, err := m.AddressList(key)

	if err != nil && err != mail.ErrHeaderNotPresent {
		return nil, err
	}
	emails := make([]*Address, len(list))

	for i, addr := range list {
		emails[i] = &Address{Name: addr.Name, Email: addr.Address}
	}
	return emails, nil
}

func readDateHeader(h mail.Header) (*Time, error) {
	date, err := h.Date()

	if err == mail.ErrHeaderNotPresent {
		return &Time{time.Now()}, nil
	}
	return &Time{date}, err
}
