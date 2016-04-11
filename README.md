# mailhook

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
```

### Payload

The following payload is posted to the web hook:

```json
{
  "mail": {
    "sender": "tina@example.io",
    "recipient": "tom+0123456789abcdef@example.com",
    "message": {
      "from": [{
        "name": "Tina",
        "email": "tina@example.io"
      }],
      "to": [{
        "name": "Tom",
        "email": "tom+0123456789abcdef@example.com"
      }, {
        "name": "Bill",
        "email": "bill@example.net"
      }],
      "cc": [{
        "name": "",
        "email": "heather@example.io"
      }],
      "bcc": [],
      "reply_to": [{
        "name": "Tina",
        "email": "tina+tag@example.net"
      }],
      "subject": "Re: This is a test message :)",
      "date": "2016-03-14T18:19:42+01:00",
      "message_id": "<1dfe5dbc-3862-4872-bb7b-470a27381d1f@example.io>",
      "in_reply_to": "<39895116-02af-47a1-b7fd-e13bf343f798@example.net>",
      "references": [
        "<6b67a383-8c72-400b-827b-b42497437f83@example.io>",
        "<39895116-02af-47a1-b7fd-e13bf343f798@example.net>"
      ],
      "text": "Hey Everybody,\n\nThis is just a test...\n\n",
      "html": "<html><body><p>Hey Everybody,</p><p>This is just a test...</p></body></html>"
    }
}
```
