package main

import (
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) MarshalMultipart() string {
	return t.Format(time.RFC3339)
}
