FROM golang:onbuild

RUN adduser --gecos GECOS --disabled-password --shell /bin/bash app

USER app

EXPOSE 1025

CMD app -listen :1025 -hook-file /etc/mailhook/hooks.yml
