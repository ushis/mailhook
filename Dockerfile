FROM alpine:3.5

RUN apk add --no-cache ca-certificates

COPY mailhook /usr/local/bin/mailhook

USER 1

EXPOSE 1025

CMD mailhook -listen :1025 -hook-dir /hooks
