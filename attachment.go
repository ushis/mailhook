package main

import (
	"github.com/jhillyerd/go.enmime"
)

type Attachment struct {
	enmime.MIMEPart
}

func (a *Attachment) MarshalMultipart() (string, []byte) {
	return a.FileName(), a.Content()
}
