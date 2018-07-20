FROM golang:1.10-alpine3.7

WORKDIR /go/src/github.com/hypnoglow/chronologist

COPY . .

RUN go build -o bin/chronologist ./cmd/chronologist

FROM alpine:3.7

RUN apk add --no-cache ca-certificates \
    && addgroup -S chronologist -g 1177 \
    && adduser -S -G chronologist chronologist -u 1177

COPY --from=0 /go/src/github.com/hypnoglow/chronologist/bin/chronologist /usr/local/bin/chronologist

USER 1177

ENTRYPOINT ["/usr/local/bin/chronologist"]

ARG VCS_REF
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/hypnoglow/chronologist" \
      org.label-schema.license="Apache-2.0" \
      org.label-schema.schema-version="1.0"
