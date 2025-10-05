FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags="-s -w" ./cmd/server/main.go


FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
