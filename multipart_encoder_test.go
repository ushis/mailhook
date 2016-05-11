package main

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	msg, err := readTestMail("attachments.eml")

	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	enc := NewMultipartEncoder(buf)
	mail := NewMail("bill@example.com", "tina@example.org", msg)

	if err := enc.Encode("mail", mail); err != nil {
		t.Fatal(err)
	}
	enc.Close()

	f, err := multipart.NewReader(buf, enc.Boundary()).ReadForm(int64(10 << 20))

	if err != nil {
		t.Fatal(err)
	}

	if n := len(f.Value["mail[sender]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[sender]"][0]; v != "bill@example.com" {
		t.Fatalf("\"bill@example.org\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[recipient]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[recipient]"][0]; v != "tina@example.org" {
		t.Fatalf("\"tina@example.org\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][from][][name]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[message][from][][name]"][0]; v != "Bob" {
		t.Fatalf("\"Bob\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][from][][email]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[message][from][][email]"][0]; v != "bob@example.org" {
		t.Fatalf("\"bob@example.org\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][to][][email]"]); n != 2 {
		t.Fatalf("2 != %d", n)
	}

	if v := f.Value["mail[message][to][][email]"][0]; v != "alice@example.com" {
		t.Fatalf("\"alice@example.com\" != \"%s\"", v)
	}

	if v := f.Value["mail[message][to][][email]"][1]; v != "tina@example.com" {
		t.Fatalf("\"tina@example.com\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][cc][][email]"]); n != 0 {
		t.Fatalf("0 != %d", n)
	}

	if n := len(f.Value["mail[message][bcc][][email]"]); n != 0 {
		t.Fatalf("0 != %d", n)
	}

	if n := len(f.Value["mail[message][reply_to][][email]"]); n != 0 {
		t.Fatalf("0 != %d", n)
	}

	if n := len(f.Value["mail[message][subject]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[message][subject]"][0]; v != "Re: This is just a test" {
		t.Fatalf("\"Re: This is just a test\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][date]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}
	loc, err := time.LoadLocation("Europe/Berlin")

	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2016, 4, 11, 17, 44, 9, 0, loc).Format(time.RFC3339)

	if v := f.Value["mail[message][date]"][0]; v != date {
		t.Fatalf("\"%s\" != \"%s\"", date, v)
	}

	if n := len(f.Value["mail[message][message_id]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[message][message_id]"][0]; v != "<5727b6c4@example.org>" {
		t.Fatalf("\"<5727b6c4@example.org>\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][in_reply_to]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[message][in_reply_to]"][0]; v != "<8b6ea071@example.com>" {
		t.Fatalf("\"<8b6ea071@example.com>\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][references][]"]); n != 3 {
		t.Fatalf("3 != %d", n)
	}

	if v := f.Value["mail[message][references][]"][0]; v != "<8ca8a3e3@example.com>" {
		t.Fatalf("\"<8ca8a3e3@example.com>\" != \"%s\"", v)
	}

	if v := f.Value["mail[message][references][]"][2]; v != "<8b6ea071@example.com>" {
		t.Fatalf("\"<8b6ea071@example.com>\" != \"%s\"", v)
	}

	if n := len(f.Value["mail[message][text]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}
	text := `Hey Alice, thanks for your test mail!

> Hey Bob,
>
> this is just a test...
>
> Cheers, Alice
`

	if v := f.Value["mail[message][text]"][0]; v != text {
		t.Fatalf("unexpected text: \"%s\"", v)
	}

	if n := len(f.Value["mail[message][html]"]); n != 1 {
		t.Fatalf("1 != %d", n)
	}

	if v := f.Value["mail[message][html]"][0]; v != "" {
		t.Fatalf("unexpected html: \"%s\"", v)
	}

	if n := len(f.File["mail[message][attachments][]"]); n != 4 {
		t.Fatalf("4 != %d", n)
	}

	for _, a := range f.File["mail[message][attachments][]"] {
		orig, err := readAttachment(a.Filename)

		if err != nil {
			t.Fatal(err)
		}
		file, err := a.Open()

		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(file)

		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(orig, body) {
			t.Fatalf("unexpected body for: %s", a.Filename)
		}
	}
}
