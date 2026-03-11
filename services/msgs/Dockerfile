FROM golang:1.24-alpine AS builder

WORKDIR /msgs

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o msgs ./services/msgs/cmd

FROM alpine:latest

WORKDIR /msgs

COPY --from=builder /msgs/msgs .

EXPOSE 8086

CMD ["./msgs"]