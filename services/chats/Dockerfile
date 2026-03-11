FROM golang:1.24-alpine AS builder

WORKDIR /chats

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chats ./services/chats/cmd

FROM alpine:latest

WORKDIR /chats

COPY --from=builder /chats/chats .

EXPOSE 8085

CMD ["./chats"]