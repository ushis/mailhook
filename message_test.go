package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func readTestMail(name string) (*Message, error) {
	f, err := os.Open(filepath.Join("testdata", name))

	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewMessage(f)
}

func TestTextMail(t *testing.T) {
	msg, err := readTestMail("text.eml")

	if err != nil {
		t.Fatal(err)
	}

	if len(msg.From) != 1 {
		t.Fatalf("1 != %d", len(msg.From))
	}

	if msg.From[0].Name != "Bob" {
		t.Fatalf("\"Bob\" != \"%s\"", msg.From[0].Name)
	}

	if msg.From[0].Email != "bob@example.org" {
		t.Fatalf("\"bob@example.org\" != \"%s\"", msg.From[0].Email)
	}

	if len(msg.To) != 1 {
		t.Fatalf("1 != %d", len(msg.To))
	}

	if msg.To[0].Name != "Alice" {
		t.Fatalf("\"Alice\" != \"%s\"", msg.To[0].Name)
	}

	if msg.To[0].Email != "alice@example.com" {
		t.Fatalf("\"alice@example.com\" != \"%s\"", msg.To[0].Email)
	}

	if len(msg.Cc) != 0 {
		t.Fatalf("0 != %d", len(msg.Cc))
	}

	if len(msg.Bcc) != 0 {
		t.Fatalf("0 != %d", len(msg.Bcc))
	}

	if len(msg.ReplyTo) != 0 {
		t.Fatalf("0 != %d", len(msg.ReplyTo))
	}

	if msg.Subject != "Re: This is just a test" {
		t.Fatalf("\"Re: This is just a test\" != \"%s\"", msg.Subject)
	}

	if msg.Date.UTC() != time.Date(2016, 4, 11, 15, 44, 9, 0, time.UTC) {
		t.Fatalf("unexpected date: %v", msg.Date)
	}

	if msg.MessageID != "<5727b6c4@example.org>" {
		t.Fatalf("\"<5727b6c4@example.org>\" != \"%s\"", msg.MessageID)
	}

	if msg.InReplyTo != "<8b6ea071@example.com>" {
		t.Fatalf("\"<8b6ea071@example.com> != \"%s\"", msg.InReplyTo)
	}

	if len(msg.References) != 3 {
		t.Fatalf("3 != %d", len(msg.References))
	}

	if msg.References[0] != "<8ca8a3e3@example.com>" {
		t.Fatalf("\"<8ca8a3e3@example.com>\" != \"%s\"", msg.References[0])
	}

	if msg.References[1] != "<c8fecb75@example.org>" {
		t.Fatalf("\"<c8fecb75@example.org>\" != \"%s\"", msg.References[1])
	}

	if msg.References[2] != "<8b6ea071@example.com>" {
		t.Fatalf("\"<8b6ea071@example.com>\" != \"%s\"", msg.References[2])
	}

	text := `Hey Alice, thanks for your test mail!

> Hey Bob,
>
> this is just a test...
>
> Cheers, Alice
`

	if msg.Text != text {
		t.Fatalf("unexpected text: \"%s\"", msg.Text)
	}

	if len(msg.HTML) != 0 {
		t.Fatalf("0 != %d", len(msg.HTML))
	}
}

func TestHTMLMail(t *testing.T) {
	msg, err := readTestMail("html.eml")

	if err != nil {
		t.Fatal(err)
	}

	text := `Hey Alice, thanks for your test mail!

Hey Bob,

this is just a test...

Cheers, Alice`

	if msg.Text != text {
		t.Fatalf("unexpected text: \"%s\"", msg.Text)
	}

	html := `<html>
  <body>
    <p>Hey Alice, thanks for your test mail!</p>

    <blockquote>
      <p>Hey Bob,</p>

      <p>this is just a test...</p>

      <p>Cheers, Alice</p>
    </blockquote>
  </body>
</html>
`

	if msg.HTML != html {
		t.Fatalf("unexpected html != \"%s\"", msg.HTML)
	}
}

func TestAlternativeMail(t *testing.T) {
	msg, err := readTestMail("alternative.eml")

	if err != nil {
		t.Fatal(err)
	}

	text := `Hey Alice, thanks for your test mail!

> Hey Bob,
>
> this is just a test...
>
> Cheers, Alice
`

	if msg.Text != text {
		t.Fatalf("unexpected text: \"%s\"", msg.Text)
	}

	html := `<html>
  <body>
    <p>Hey Alice, thanks for your test mail!</p>

    <blockquote>
      <p>Hey Bob,</p>

      <p>this is just a test...</p>

      <p>Cheers, Alice</p>
    </blockquote>
  </body>
</html>
`

	if msg.HTML != html {
		t.Fatalf("unexpected html != \"%s\"", msg.HTML)
	}
}
