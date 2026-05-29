# syntax=docker/dockerfile:1.4
FROM golang:1.26

ARG TARGETPLATFORM

ENV GOROOT /usr/local/go

ENV GOTOOLCHAIN auto

RUN git config --global --add safe.directory '*'

COPY $TARGETPLATFORM/commitlint-scope /usr/bin/
CMD ["commitlint-scope"]
