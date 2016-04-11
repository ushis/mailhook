package main

import (
  "time"
  "net/mail"
  "io"
  "github.com/jhillyerd/go.enmime"
  "strings"
)

type Message struct {
	From       []string `json:"from"`
	To         []string `json:"to"`
  Cc         []string `json:"cc"`
  Bcc        []string `json:"bcc"`
  ReplyTo    []string `json:"reply_to"`
	Subject    string   `json:"subject"`
	Date       time.Time   `json:"date"`
	MessageID  string   `json:"message_id"`
	InReplyTo  string   `json:"in_reply_to"`
	References []string `json:"references"`
  Text       string   `json:"text"`
  HTML       string   `json:"html"`
}

func NewMessage(r io.Reader) (*Message, error) {
  msg, err := mail.ReadMessage(r)

  if err != nil {
    return nil, err
  }
  body, err := enmime.ParseMIMEBody(msg)

  if err != nil {
    return nil, err
  }

  m := &Message{
    Subject: body.GetHeader("Subject"),
    MessageID: body.GetHeader("Message-ID"),
    InReplyTo: body.GetHeader("In-Reply-To"),
    References: strings.Fields(body.GetHeader("References")),
    Text: body.Text,
    HTML: body.HTML,
  }

  if m.From, err = readAddressListHeader(body, "From"); err != nil {
    return nil, err
  }

  if m.To, err = readAddressListHeader(body, "To"); err != nil {
    return nil, err
  }

  if m.Cc, err = readAddressListHeader(body, "Cc"); err != nil {
    return nil, err
  }

  if m.Bcc, err = readAddressListHeader(body, "Bcc"); err != nil {
    return nil, err
  }

  if m.ReplyTo, err = readAddressListHeader(body, "Reply-To"); err != nil {
    return nil, err
  }

  if m.Date, err = readDateHeader(msg.Header); err != nil {
    return nil, err
  }
  return m, nil
}

func readAddressListHeader(m *enmime.MIMEBody, key string) ([]string, error) {
  list, err := m.AddressList(key)

  if err == mail.ErrHeaderNotPresent {
    return []string{}, nil
  }

  if err != nil {
    return nil, err
  }
  addrs := make([]string, len(list))

  for i, addr := range list {
    addrs[i] = addr.String()
  }
  return addrs, nil
}

func readDateHeader(h mail.Header) (time.Time, error) {
  date, err := h.Date()

  if err == mail.ErrHeaderNotPresent {
    return time.Now(), nil
  }
  return date, err
}
