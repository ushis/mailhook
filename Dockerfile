FROM golang:onbuild

USER 1

EXPOSE 1025

CMD app -listen :1025 -hook-file /etc/mailhook/hooks.yml
