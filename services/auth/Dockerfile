FROM golang:1.24-alpine AS builder

WORKDIR /auth

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth ./services/auth/cmd

FROM alpine:latest

WORKDIR /auth

COPY --from=builder /auth/auth .

EXPOSE 8087

CMD ["./auth"]