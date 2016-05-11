# mailhook

[![Circle CI](https://circleci.com/gh/ushis/mailhook.svg?style=svg)](https://circleci.com/gh/ushis/mailhook)

Micro service receiving emails and posting them to web hooks.

## Usage

```
mailhook -hook-file hooks.yml -listen :25
```

### Hook File

A hook file includes a YAML encoded set of web hooks and email address
patterns.

```yaml
---
# catch mails to tom@example.com and tom+arbitrary-tag@example.com
- hook: 'https://api.example.com/v1/mails'
  emails:
    - 'tom@example.com'

# catch mails to *@example.net and *@example.org
- hook: 'http://example.net/messages'
  emails:
    - '@example.net'
    - '@example.org'

# catch all mails
- hook: 'http://example.net/fallback'
  emails:
    - '@'
```

### Payload

The request that will be sent to the HTTP endpoint will be encoded as
multipart/form-data with the following payload:

| Field | Value |
| ----- | ----- |
| ```mail[sender]``` | SMTP sender address |
| ```mail[recipient]``` | SMTP recipient address |
| ```mail[message][from][][name]``` | Names taken from the From header |
| ```mail[message][from][][email]``` | Adresses takem from the From header |
| ```mail[message][to][][name]``` | Names taken from the To header |
| ```mail[message][to][][email]``` | Adresses takem from the To header |
| ```mail[message][cc][][name]``` | Names taken from the Cc header |
| ```mail[message][cc][][email]``` | Adresses takem from the Cc header |
| ```mail[message][bcc][][name]``` | Names taken from the Bcc header |
| ```mail[message][bcc][][email]``` | Adresses takem from the Bcc header |
| ```mail[message][reply_to][][name]``` | Names taken from the ReplyTo header |
| ```mail[message][reply_to][][email]``` | Adresses takem from the ReplyTo header |
| ```mail[message][subject]``` | Subject of the message |
| ```mail[message][date]``` | RFC3339 encoded date of the message |
| ```mail[message][message_id]``` | ID taken from the MessageID header |
| ```mail[message][in_reply_to]``` | ID taken from the InReplyTo header |
| ```mail[message][references][]``` | IDs taken from the References header |
| ```mail[message][text]``` | Text body of the message |
| ```mail[message][html]``` | HTML body of the message (can be empty) |
| ```mail[message][attachments][]``` | Attachments of the message |
