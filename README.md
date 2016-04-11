# mailhook

## Usage

```
mailhook -hook-file hooks.yml -listen :25
```

### Hook File

```yaml
---
- hook: 'https://api.example.com/v1/mails'
  emails:
    # Catches:
    #   - tom@example.com
    #   - tom+arbitrary-tag@example.com
    - 'tom@example.com'

- hook: 'http://example.net/messages'
  emails:
    # Catches:
    #   - *@example.net
    #   - *@example.org
    - @example.net
    - @example.org
```

### Payload

The following payload is posted to the web hook:

```json
{
  "mail": {
    "from": [
      "<tina@example.io>"
    ],
    "to": [
      "<tom+0123456789abcdef@example.com>",
      "<bill@example.net>"
    ],
    "cc": [
      "<heather@example.io"
    ],
    "bcc": [],
    "reply_to": [
      "<tina+tag@example.net>"
    ],
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
