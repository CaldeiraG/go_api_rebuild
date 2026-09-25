FROM golang:1.22 as builder

RUN apt update && \
    apt install -y --no-install-recommends \
    ca-certificates

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o go-api-template ./main.go
RUN mv go-api-template go-api-template.bin

FROM ubuntu:22.04

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/go-api-template.bin /srv/go-api-template
COPY --from=builder /app/static /srv/static
COPY --from=builder /app/prodFAssy_config.json /srv/prodFAssy_config.json

# Run from /srv so relative paths (./static, prodFAssy_config.json) resolve.
WORKDIR /srv

EXPOSE 4000

CMD ["/srv/go-api-template"]
