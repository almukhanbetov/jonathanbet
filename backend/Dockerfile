FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
EXPOSE 8282
