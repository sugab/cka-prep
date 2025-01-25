# Builder stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# Runtime image stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
CMD [ "./server" ]
