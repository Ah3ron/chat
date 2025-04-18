# --- builder stage ---
FROM golang:1.20-alpine AS builder
WORKDIR /app
# fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o chat-server ./main.go

# --- final stage ---
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/chat-server .
EXPOSE 3000
CMD ["./chat-server"]