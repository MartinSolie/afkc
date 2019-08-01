FROM golang:1.12 AS builder

WORKDIR /go/src/github.com/martinsolie/afkc
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -o /usr/local/bin/afkc github.com/martinsolie/afkc

# ===

FROM alpine:3.9

COPY --from=builder /usr/local/bin/afkc /usr/local/bin/afkc
RUN apk add --no-cache ca-certificates

EXPOSE 8000

ENTRYPOINT ["afkc"]
