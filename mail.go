package main

type Mail struct {
	Sender    string   `multipart:"sender"`
	Recipient string   `multipart:"recipient"`
	Message   *Message `multipart:"message"`
}

func NewMail(sender, recipient string, msg *Message) *Mail {
	return &Mail{Sender: sender, Recipient: recipient, Message: msg}
}
