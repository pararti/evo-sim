# Stage 1: Build
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o evosim cmd/app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/evosim .
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./evosim"]