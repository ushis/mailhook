package main

type Payload struct {
	Mail *Mail `json:"mail"`
}

type Mail struct {
	Sender    string   `json:"sender"`
	Recipient string   `json:"recipient"`
	Message   *Message `json:"message"`
}

func NewPayload(sender, recipient string, msg *Message) *Payload {
	return &Payload{
		Mail: &Mail{
			Sender:    sender,
			Recipient: recipient,
			Message:   msg,
		},
	}
}
