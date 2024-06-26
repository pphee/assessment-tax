# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN gofmt -w .
RUN go vet ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]

