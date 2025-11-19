FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/gateway

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8283
ENTRYPOINT ["./app"]