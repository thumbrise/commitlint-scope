# syntax=docker/dockerfile:1.4
FROM alpine:3.21

ARG TARGETPLATFORM

RUN apk --no-cache add git

RUN git config --global --add safe.directory '*'

COPY $TARGETPLATFORM/commitlint-scope /usr/bin/

ENTRYPOINT ["commitlint-scope"]
CMD []
