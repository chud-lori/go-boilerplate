# Start with an official Golang image
FROM golang:1.22.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

# Final stage for a smaller production image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
