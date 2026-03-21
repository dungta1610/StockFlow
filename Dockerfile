# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for fetching private/public Go modules if needed
RUN apk add --no-cache git

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stockflow main.go

# Runtime stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/stockflow /app/stockflow
COPY .env /app/.env

EXPOSE 8080

CMD ["/app/stockflow"]