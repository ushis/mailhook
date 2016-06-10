FROM alpine

RUN apk update && apk add ca-certificates

COPY mailhook /usr/local/bin/mailhook

USER 1

EXPOSE 1025

CMD mailhook -listen :1025 -hook-file /etc/mailhook/hooks.yml
