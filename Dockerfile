FROM golang:1.16 AS builder

WORKDIR /go/src/github.com/whywaita/ingress-external-caddy

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY . .
RUN go build -o ./ingress-external-caddy .

FROM alpine

RUN apk update \
  && apk update
RUN apk add --no-cache ca-certificates \
  && update-ca-certificates 2>/dev/null || true

COPY --from=builder /go/src/github.com/whywaita/ingress-external-caddy/ingress-external-caddy /app

CMD ["/app"]