FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go mod download
RUN go build -o forum ./cmd/web
FROM alpine:3.6
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
CMD ["./forum"]

# CMD ["go", "run", "./cmd/web"]

