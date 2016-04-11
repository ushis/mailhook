FROM golang:onbuild

EXPOSE 25

CMD app -listen :25 -hook-file /hooks.yml
