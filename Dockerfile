FROM golang:1.16 as builder
COPY . /tmp/gostaticwebserver
WORKDIR /tmp/gostaticwebserver

RUN CGO_ENABLED=0 GOOS=linux go build -a -o webserver cmd/gostaticwebserver/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl jq bash emacs-nox

WORKDIR /webserver/
COPY --from=builder /tmp/gostaticwebserver/webserver .

CMD ["./webserver"]
