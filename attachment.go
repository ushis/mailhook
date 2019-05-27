package main

import (
	"github.com/jhillyerd/enmime"
)

type Attachment struct {
	*enmime.Part
}

func (a *Attachment) MarshalMultipart() (string, []byte) {
	return a.FileName, a.Content
}
