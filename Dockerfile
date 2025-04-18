FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
RUN go mod download && go build -o server cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .

USER 1001

EXPOSE 11434
CMD ["/app/server"]
