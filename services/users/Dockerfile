FROM golang:1.24-alpine AS builder

WORKDIR /users

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o users ./services/users/cmd

FROM alpine:latest

WORKDIR /users

COPY --from=builder /users/users .

EXPOSE 8082

CMD ["./users"]