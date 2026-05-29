# syntax=docker/dockerfile:1.4
FROM golang:1.26-alpine

ARG TARGETPLATFORM

ENV GOROOT /usr/local/go

ENV GOTOOLCHAIN auto

RUN apk --no-cache add gcc musl-dev git mercurial

RUN git config --global --add safe.directory '*'

COPY $TARGETPLATFORM/commitlint-scope /usr/bin/
CMD ["commitlint-scope"]
