FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/msg-processor

FROM alpine
WORKDIR /app
COPY --from=builder /app/msg-processor /app/msg-processor
EXPOSE 8080
CMD ["/app/msg-processor"]

