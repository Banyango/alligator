FROM golang:1.10.1 AS builder
RUN go version

COPY . /go/src/github.com/Banyango/Alligator
WORKDIR /go/src/github.com/Banyango/Alligator
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep init && \
    dep ensure -v

RUN make build-linux

FROM alpine:latest
COPY --from=builder /go/src/github.com/Banyango/Alligator/.dist/ .
ENV GODEBUG netdns=go
EXPOSE 8080
ENTRYPOINT ["./alligator"]